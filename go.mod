module github.com/rancher/k3os

go 1.12

require (
	github.com/docker/docker v1.13.1
	github.com/ghodss/yaml v1.0.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mattn/go-isatty v0.0.7
	github.com/pkg/errors v0.8.1
	github.com/rancher/mapper v0.0.0-20190417000032-48be2d1eadc0
	github.com/rancher/wrangler v0.2.0
	github.com/rancher/wrangler-api v0.2.0
	github.com/sclevine/spec v1.3.0
	github.com/sirupsen/logrus v1.4.1
	github.com/urfave/cli v1.20.0
	golang.org/x/crypto v0.0.0-20190422183909-d864b10871cd
	golang.org/x/sys v0.0.0-20190418153312-f0ce4c0180be
	gopkg.in/freddierice/go-losetup.v1 v1.0.0-20170407175016-fc9adea44124
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)
