package agent

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	v1 "github.com/rancher/k3os/pkg/apis/k3os.cattle.io/v1"
	"github.com/rancher/k3os/pkg/controller"
	"github.com/rancher/k3os/pkg/generated/controllers/k3os.cattle.io"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/batch"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/core"
	"github.com/rancher/wrangler/pkg/resolvehome"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/rancher/wrangler/pkg/start"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// Command is the `agent` sub-command, it is the k3OS resource controller.
var Command = cli.Command{
	Name:  "agent",
	Usage: "control custom resources",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:   "threads",
			EnvVar: "THREADS",
			Value:  1,
		},
		cli.StringFlag{
			Name:   "kubeconfig",
			EnvVar: "KUBECONFIG",
		},
		cli.StringFlag{
			Name:   "namespace",
			EnvVar: "NAMESPACE",
		},
		cli.StringFlag{
			Name:   "masterurl",
			EnvVar: "MASTERURL",
			Value:  "",
		},
	},
	Before: func(c *cli.Context) error {
		// required parameters
		if ns := c.String("namespace"); len(ns) == 0 {
			return errors.New("namespace is required")
		}
		// required uid
		if os.Getuid() != 0 {
			return errors.New("must be run as root")
		}
		// required filesystem
		systemDir := "/k3os/system"
		if inf, err := os.Stat(systemDir); err != nil {
			return err
		} else if !inf.IsDir() {
			return fmt.Errorf("stat %s: not a directory", systemDir)
		}
		return nil
	},
	Action: Run,
}

// Run the `agent` sub-command
func Run(c *cli.Context) {
	logrus.Debug("K3OS::OPERATOR >>> SETUP")

	// ensure that we are running k3os
	current, err := os.Readlink("/k3os/system/k3os/current")
	if err != nil {
		logrus.Fatal(err)
	}
	current = filepath.Base(current)

	kubeconfig, err := resolvehome.Resolve(c.String("kubeconfig"))
	if err != nil {
		logrus.Info("Resolving home dir failed.")
	}
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		kubeconfig = ""
	}

	threads := c.Int("threads")
	masterurl := c.String("masterurl")
	namespace := c.String("namespace")

	ctx := signals.SetupSignalHandler(context.Background())

	cfg, err := clientcmd.BuildConfigFromFlags(masterurl, kubeconfig)
	if err != nil {
		logrus.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	k3osFactory, err := k3os.NewFactoryFromConfigWithNamespace(cfg, namespace)
	if err != nil {
		logrus.Fatalf("Error building k3OS controllers: %s", err.Error())
	}

	coreFactory, err := core.NewFactoryFromConfigWithNamespace(cfg, namespace)
	if err != nil {
		logrus.Fatalf("Error building core controllers: %s", err.Error())
	}

	batchFactory, err := batch.NewFactoryFromConfigWithNamespace(cfg, namespace)
	if err != nil {
		logrus.Fatalf("Error building rbac controllers: %s", err.Error())
	}

	logrus.Debug("K3OS::OPERATOR >>> REGISTER")

	updateChannelController := k3osFactory.K3os().V1().UpdateChannel()
	controller.Register(ctx,
		updateChannelController,
		coreFactory.Core().V1().Node(),
		batchFactory.Batch().V1().Job(),
	)

	if err := start.All(ctx, threads, k3osFactory, coreFactory, batchFactory); err != nil {
		logrus.Fatalf("Error starting: %s", err.Error())
	}

	if list, err := updateChannelController.List(namespace, metav1.ListOptions{Limit: 1}); err != nil {
		logrus.Warn(err)
	} else if len(list.Items) == 0 {
		if upchan, err := updateChannelController.Create(&v1.UpdateChannel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "github-releases",
				Namespace: namespace,
				Annotations: map[string]string{
					"k3os.io/node": os.Getenv("K3OS_NODE_NAME"),
				},
			},
			Spec: v1.UpdateChannelSpec{
				URL:         "github-releases://rancher/k3os",
				Concurrency: 1,
				Version:     current,
			},
		}); err != nil {
			logrus.Warn(err)
		} else {
			logrus.Infof("Created default UpdateChannel: name=%s, url=%s, version=%s, concurrency=%d",
				upchan.ObjectMeta.Name,
				upchan.Spec.URL,
				upchan.Spec.Version,
				upchan.Spec.Concurrency,
			)
		}
	}

	<-ctx.Done()
}
