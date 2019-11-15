package main

import (
	v1 "github.com/rancher/k3os/pkg/apis/k3os.cattle.io/v1"
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	controllergen.Run(args.Options{
		OutputPackage: "github.com/rancher/k3os/pkg/generated",
		Boilerplate:   "hack/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"k3os.cattle.io": {
				Types: []interface{}{
					v1.UpdateChannel{},
				},
				GenerateTypes: true,
			},
		},
	})
}
