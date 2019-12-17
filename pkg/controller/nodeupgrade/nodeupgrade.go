package nodeupgrade

import (
	"context"
	"fmt"

	api "github.com/rancher/k3os/pkg/apis/k3os.io"
	apiv1 "github.com/rancher/k3os/pkg/apis/k3os.io/v1"
	ctl "github.com/rancher/k3os/pkg/generated/controllers/k3os.io"
	ctlv1 "github.com/rancher/k3os/pkg/generated/controllers/k3os.io/v1"
	"github.com/rancher/norman/pkg/openapi"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/batch"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/core"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

type Options struct {
	ControllerName     string
	ControllerApply    apply.Apply
	ControllerFactory  *ctl.Factory
	CoreFactory        *core.Factory
	BatchFactory       *batch.Factory
	ServiceAccountName string
}

// BatchCreateCRD registers the Channel CRD
func BatchCreateCRD(ctx context.Context, factory *crd.Factory) error {
	prototype := apiv1.NewNodeUpgrade("", "", apiv1.NodeUpgrade{})
	schema, err := openapi.ToOpenAPIFromStruct(*prototype)
	if err != nil {
		return err
	}
	factory.BatchCreateCRDs(ctx, crd.CRD{
		GVK:        prototype.GroupVersionKind(),
		PluralName: apiv1.NodeUpgradeResourceName,
		Status:     true,
		Schema:     schema,
		Categories: []string{"all", "k3os", "upgrade"},
		ShortNames: []string{"noup", "noups"},
		//Columns: []v1beta1.CustomResourceColumnDefinition{
		//	{
		//		Name:     "Node",
		//		Type:     "string",
		//		JSONPath: ".spec.nodeName",
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

// RegisterHandlers registers NodeUpgrade handlers
func RegisterHandlers(ctx context.Context, options Options) {
	jobs := options.BatchFactory.Batch().V1().Job()
	nodeUpgrades := options.ControllerFactory.K3os().V1().NodeUpgrade()

	ctlv1.RegisterNodeUpgradeStatusHandler(ctx, nodeUpgrades, "", options.ControllerName,
		statusHandler(ctx, options),
	)

	ctlv1.RegisterNodeUpgradeGeneratingHandler(ctx, nodeUpgrades, options.ControllerApply.WithCacheTypes(jobs), "", options.ControllerName,
		generatingHandler(ctx, options), &generic.GeneratingHandlerOptions{},
	)
}

func statusHandler(ctx context.Context, options Options) ctlv1.NodeUpgradeStatusHandler {
	return func(obj *apiv1.NodeUpgrade, status apiv1.NodeUpgradeStatus) (apiv1.NodeUpgradeStatus, error) {
		if apiv1.NodeUpgradeScheduled.IsTrue(obj) {
			nodeController := options.CoreFactory.Core().V1().Node()
			node, err := nodeController.Get(obj.Spec.NodeName, metav1.GetOptions{})
			if err != nil {
				logrus.Error(err)
			}
			node.Labels = labels.Merge(node.Labels, labels.Set{
				api.LabelUpgradeVersion: obj.Spec.Version,
			})
			node, err = nodeController.Update(node)
			logrus.Debugf("%#v", node)
			return status, err
		}
		return status, nil
	}
}

func generatingHandler(ctx context.Context, options Options) ctlv1.NodeUpgradeGeneratingHandler {
	return func(obj *apiv1.NodeUpgrade, status apiv1.NodeUpgradeStatus) (objects []runtime.Object, _ apiv1.NodeUpgradeStatus, _ error) {
		apiv1.NodeUpgradeScheduled.CreateUnknownIfNotExists(obj)
		if apiv1.NodeUpgradeScheduled.IsUnknown(obj) {
			apiv1.NodeUpgradeScheduled.True(obj)
		}
		job := newUpgradeJob(obj.Spec, options)
		job.Name = fmt.Sprintf("%s-job", obj.Name)
		return append(objects, job), obj.Status, nil
	}
}
