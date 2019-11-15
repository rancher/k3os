package controller

import (
	"context"

	"github.com/rancher/k3os/pkg/controller/updatechannel"
	k3osv1 "github.com/rancher/k3os/pkg/generated/controllers/k3os.cattle.io/v1"
	batchv1 "github.com/rancher/wrangler-api/pkg/generated/controllers/batch/v1"
	corev1 "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"
)

func Register(
	ctx context.Context,
	updateChannels k3osv1.UpdateChannelController,
	nodes corev1.NodeController,
	jobs batchv1.JobController,
) {
	handler := updatechannel.NewHandler(ctx, updateChannels, nodes, jobs)
	updateChannels.OnChange(ctx, "k3os-operator", handler.OnChange)
	updateChannels.OnRemove(ctx, "k3os-operator", handler.OnRemove)
	jobs.OnChange(ctx, "k3os-operator", handler.JobOnChange)
	jobs.OnRemove(ctx, "k3os-operator", handler.JobOnChange)
}
