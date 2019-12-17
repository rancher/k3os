package controller

import (
	"context"

	"github.com/rancher/k3os/pkg/controller/channel"
	"github.com/rancher/k3os/pkg/controller/nodeupgrade"
	"github.com/rancher/k3os/pkg/controller/upgradeset"
	"github.com/rancher/k3os/pkg/generated/controllers/k3os.io"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/batch"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/core"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/start"
	"k8s.io/client-go/rest"
)

type Options struct {
	Name               string
	Namespace          string
	ServiceAccountName string
	Threads            int
}

// Start the controller
func Start(ctx context.Context, options Options) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	if err = registerCRDs(ctx, config); err != nil {
		return err
	}

	coreFactory, err := core.NewFactoryFromConfigWithNamespace(config, options.Namespace)
	if err != nil {
		return err
	}
	batchFactory, err := batch.NewFactoryFromConfigWithNamespace(config, options.Namespace)
	if err != nil {
		return err
	}
	factory, err := k3os.NewFactoryFromConfigWithNamespace(config, options.Namespace)
	if err != nil {
		return err
	}
	apply, err := apply.NewForConfig(config)
	if err != nil {
		return err
	}

	channel.RegisterHandlers(ctx, channel.Options{
		ControllerName:    options.Name,
		ControllerFactory: factory,
		ControllerApply:   apply,
		PollingInterval:   channel.DefaultPollingInterval,
	})

	upgradeset.RegisterHandlers(ctx, upgradeset.Options{
		ControllerName:    options.Name,
		ControllerFactory: factory,
		ControllerApply:   apply,
		CoreFactory:       coreFactory,
		ResyncInterval:    upgradeset.DefaultResyncInterval,
	})

	nodeupgrade.RegisterHandlers(ctx, nodeupgrade.Options{
		ControllerName:     options.Name,
		ControllerFactory:  factory,
		ControllerApply:    apply,
		CoreFactory:        coreFactory,
		BatchFactory:       batchFactory,
		ServiceAccountName: options.ServiceAccountName,
	})

	return start.All(ctx, options.Threads, coreFactory, batchFactory, factory)
}

func registerCRDs(ctx context.Context, config *rest.Config) error {
	factory, err := crd.NewFactoryFromClient(config)
	if err != nil {
		return err
	}
	if err = channel.BatchCreateCRD(ctx, factory); err != nil {
		return err
	}
	if err = upgradeset.BatchCreateCRD(ctx, factory); err != nil {
		return err
	}
	if err = nodeupgrade.BatchCreateCRD(ctx, factory); err != nil {
		return err
	}
	return factory.BatchWait()
}
