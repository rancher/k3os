package main

import (
	v1 "github.com/rancher/k3os/pkg/apis/k3os.io/v1"
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	controllergen.Run(args.Options{
		OutputPackage: "github.com/rancher/k3os/pkg/generated",
		Boilerplate:   "hack/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"k3os.io": {
				Types: []interface{}{
					v1.Channel{},
					v1.UpgradeSet{},
					v1.NodeUpgrade{},
				},
				GenerateTypes: true,
			},
		},
	})
}
