package channel

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	api "github.com/rancher/k3os/pkg/apis/k3os.io"
	apiv1 "github.com/rancher/k3os/pkg/apis/k3os.io/v1"
	"github.com/rancher/k3os/pkg/controller/upgradeset"
	ctl "github.com/rancher/k3os/pkg/generated/controllers/k3os.io"
	ctlv1 "github.com/rancher/k3os/pkg/generated/controllers/k3os.io/v1"
	"github.com/rancher/norman/pkg/openapi"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	// DefaultPollingInterval is the default duration between attempts to resolve latest.
	DefaultPollingInterval = 5 * time.Minute // TODO review default polling interval for channels
)

// BatchCreateCRD registers the Channel CRD
func BatchCreateCRD(ctx context.Context, factory *crd.Factory) error {
	prototype := apiv1.NewChannel("", "", apiv1.Channel{})
	schema, err := openapi.ToOpenAPIFromStruct(*prototype)
	if err != nil {
		return err
	}
	factory.BatchCreateCRDs(ctx, crd.CRD{
		GVK:        prototype.GroupVersionKind(),
		PluralName: apiv1.ChannelResourceName,
		Status:     true,
		Schema:     schema,
		Categories: []string{"all", "k3os", "upgrade"},
		ShortNames: []string{"ch", "chs", "chan", "chans"},
		//Columns: []v1beta1.CustomResourceColumnDefinition{
		//	{
		//		Name:     "Latest",
		//		Type:     "string",
		//		JSONPath: ".status.latestVersion",
		//		Priority: 10,
		//	},
		//	{
		//		Name:     "URL",
		//		Type:     "string",
		//		JSONPath: ".spec.url",
		//		Priority: 100,
		//	},
		//},
	})
	return nil
}

type Options struct {
	ControllerName    string
	ControllerApply   apply.Apply
	ControllerFactory *ctl.Factory
	PollingInterval   time.Duration
}

// RegisterHandlers registers Channel handlers
func RegisterHandlers(ctx context.Context, options Options) {
	channels := options.ControllerFactory.K3os().V1().Channel()
	upgradeSets := options.ControllerFactory.K3os().V1().UpgradeSet()

	ctlv1.RegisterChannelStatusHandler(ctx, channels, apiv1.ChannelLatestResolved, options.ControllerName,
		statusHandler(ctx, options),
	)
	ctlv1.RegisterChannelGeneratingHandler(ctx, channels, options.ControllerApply.WithCacheTypes(upgradeSets), "", options.ControllerName,
		generatingHandler(ctx, options), &generic.GeneratingHandlerOptions{},
	)
}

func generatingHandler(ctx context.Context, options Options) ctlv1.ChannelGeneratingHandler {
	timeout := upgradeset.DefaultDrainTimeout
	return func(obj *apiv1.Channel, status apiv1.ChannelStatus) (objects []runtime.Object, _ apiv1.ChannelStatus, _ error) {
		if obj.Status.LatestVersion == "" {
			return objects, status, fmt.Errorf("unresolved channel")
		}
		objects = []runtime.Object{apiv1.NewUpgradeSet(obj.Namespace, "upgrade-"+obj.Name, apiv1.UpgradeSet{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels.Set{
					api.LabelUpgradeChannel: obj.Name,
				},
			},
			Spec: apiv1.UpgradeSetSpec{
				Version:     status.LatestVersion,
				Concurrency: upgradeset.DefaultConcurrency,
				Drain: &apiv1.DrainSpec{
					Force:            true,
					Timeout:          &timeout,
					IgnoreDaemonSets: true,
					DeleteLocalData:  true,
				},
			},
		})}
		logrus.Debugf("%#v", objects)
		return objects, status, nil
	}
}

func statusHandler(ctx context.Context, options Options) ctlv1.ChannelStatusHandler {
	httpClient := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	channels := options.ControllerFactory.K3os().V1().Channel()
	return func(obj *apiv1.Channel, status apiv1.ChannelStatus) (apiv1.ChannelStatus, error) {
		if obj.Spec.URL == "" {
			return status, fmt.Errorf("missing url")
		}
		logrus.Debugf("Resolving %v", obj.Spec.URL)

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, obj.Spec.URL, nil)
		if err != nil {
			return status, err
		}
		sysid := obj.GetClusterName()
		if len(sysid) > 0 {
			sysid = fmt.Sprintf("cluster:%s", sysid)
		} else {
			sysid = fmt.Sprintf("channel:%v", obj.GetUID())
		}
		request.Header.Set("x-api-system", sysid)

		response, err := httpClient.Do(request)
		if err != nil {
			return status, err
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusFound {
			loc, err := response.Location()
			if err != nil {
				return status, err
			}
			status.LatestVersion = filepath.Base(loc.Path)
		} else {
			status.LatestVersion = filepath.Base(obj.Spec.URL)
		}

		channels.EnqueueAfter(obj.Namespace, obj.Name, options.PollingInterval)
		logrus.Debugf("%#v", status)
		return status, nil
	}
}
