package upgradeset

import (
	"context"
	"fmt"
	"sort"
	"time"

	api "github.com/rancher/k3os/pkg/apis/k3os.io"
	apiv1 "github.com/rancher/k3os/pkg/apis/k3os.io/v1"
	ctl "github.com/rancher/k3os/pkg/generated/controllers/k3os.io"
	ctlv1 "github.com/rancher/k3os/pkg/generated/controllers/k3os.io/v1"
	"github.com/rancher/norman/pkg/openapi"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/core"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	DefaultConcurrency    = uint64(1)
	DefaultDrainTimeout   = 300 * time.Second
	DefaultResyncInterval = 30 * time.Second
)

type Options struct {
	ControllerName    string
	ControllerApply   apply.Apply
	ControllerFactory *ctl.Factory
	CoreFactory       *core.Factory
	ResyncInterval    time.Duration
}

// BatchCreateCRD registers the Channel CRD
func BatchCreateCRD(ctx context.Context, factory *crd.Factory) error {
	prototype := apiv1.NewUpgradeSet("", "", apiv1.UpgradeSet{})
	schema, err := openapi.ToOpenAPIFromStruct(*prototype)
	if err != nil {
		return err
	}
	factory.BatchCreateCRDs(ctx, crd.CRD{
		GVK:        prototype.GroupVersionKind(),
		PluralName: apiv1.UpgradeSetResourceName,
		Status:     true,
		Schema:     schema,
		Categories: []string{"all", "k3os", "upgrade"},
		ShortNames: []string{"ug", "ugs", "upset", "upsets"},
		//Columns: []v1beta1.CustomResourceColumnDefinition{
		//	{
		//		Name:     "Concurrency",
		//		Type:     "integer",
		//		JSONPath: ".spec.concurrency",
		//		Format:   "int64",
		//		Priority: 10,
		//	},
		//	{
		//		Name:     "Version",
		//		Type:     "string",
		//		JSONPath: ".spec.version",
		//		Priority: 10,
		//	},
		//},
	})
	return nil
}

// RegisterHandlers registers UpgradeSet handlers
func RegisterHandlers(ctx context.Context, options Options) {
	upgradeSets := options.ControllerFactory.K3os().V1().UpgradeSet()
	nodeUpgrades := options.ControllerFactory.K3os().V1().NodeUpgrade()
	nodeUpgradesApply := options.ControllerApply.WithCacheTypes(nodeUpgrades)

	ctlv1.RegisterUpgradeSetGeneratingHandler(ctx, upgradeSets, nodeUpgradesApply, "", options.ControllerName,
		generatingHandler(ctx, options), &generic.GeneratingHandlerOptions{
			AllowClusterScoped: true,
		},
	)
}

// generatingHandler returns an array of NodeUpgrades based on the triggering UpgradeSet and matching Nodes
func generatingHandler(ctx context.Context, options Options) ctlv1.UpgradeSetGeneratingHandler {
	nodeCache := options.CoreFactory.Core().V1().Node().Cache()
	upgradeSets := options.ControllerFactory.K3os().V1().UpgradeSet()

	return func(obj *apiv1.UpgradeSet, status apiv1.UpgradeSetStatus) (objects []runtime.Object, _ apiv1.UpgradeSetStatus, _ error) {
		var (
			nodeNames []string
		)

		nodeSelector, err := nodeSelector(obj.Spec.Version)
		if err != nil {
			logrus.Error(err)
			return objects, status, nil
		}

		if len(status.Upgrades) > 0 {
			requirementUpgrading, err := labels.NewRequirement(corev1.LabelHostname, selection.In, status.Upgrades)
			if err != nil {
				logrus.Error(err)
				return objects, status, nil
			}

			if upgradingNodes, err := nodeCache.List(nodeSelector.Add(*requirementUpgrading)); err != nil {
				logrus.Error(err)
			} else {
				logrus.Debugf("upgradingNodes = %#v", upgradingNodes)
				for _, node := range upgradingNodes {
					//nodes[node.Name] = node
					nodeNames = append(nodeNames, node.Name)
				}
			}

			requirementNotUpgrading, err := labels.NewRequirement(corev1.LabelHostname, selection.NotIn, status.Upgrades)
			if err != nil {
				logrus.Error(err)
				return objects, status, nil
			}
			nodeSelector = nodeSelector.Add(*requirementNotUpgrading)
		}

		if candidateNodes, err := nodeCache.List(nodeSelector); err != nil {
			logrus.Error(err)
		} else {
			logrus.Debugf("candidateNodes = %#v", candidateNodes)
			for i := 0; i < len(candidateNodes) && uint64(len(nodeNames)) < obj.Spec.Concurrency; i++ {
				nodeNames = append(nodeNames, candidateNodes[i].Name)
			}
		}
		logrus.Debugf("nodeNames = %q", nodeNames)

		sort.Strings(nodeNames)
		for _, nodeName := range nodeNames {
			nodeUpgrade := apiv1.NewNodeUpgrade(obj.Namespace, fmt.Sprintf("%s-%.16s", obj.Name, nodeName), apiv1.NodeUpgrade{
				ObjectMeta: func() metav1.ObjectMeta {
					objectMeta := metav1.ObjectMeta{
						Labels: labels.Set{},
					}
					if channelName, ok := obj.Labels[api.LabelUpgradeChannel]; ok {
						objectMeta.Labels[api.LabelUpgradeChannel] = channelName
					}
					return objectMeta
				}(),
				Spec: apiv1.NodeUpgradeSpec{
					NodeName: nodeName,
					Version:  obj.Spec.Version,
					Drain:    obj.Spec.Drain,
				},
			})
			objects = append(objects, nodeUpgrade)
		}
		obj.Status.Upgrades = nodeNames
		upgradeSets.EnqueueAfter(obj.Namespace, obj.Name, options.ResyncInterval)
		logrus.Debugf("%#v", objects)
		return objects, obj.Status, nil
	}
}

func nodeSelector(version string) (labels.Selector, error) {
	modeNotLive, err := labels.NewRequirement(api.LabelMode, selection.NotIn, []string{"live"})
	if err != nil {
		return nil, err
	}
	upgradeEnabled, err := labels.NewRequirement(api.LabelUpgradeEnabled, selection.In, []string{"true"})
	if err != nil {
		return nil, err
	}
	versionMismatch, err := labels.NewRequirement(api.LabelVersion, selection.NotIn, []string{version})
	if err != nil {
		return nil, err
	}
	return labels.NewSelector().Add(*modeNotLive).Add(*upgradeEnabled).Add(*versionMismatch), nil
}
