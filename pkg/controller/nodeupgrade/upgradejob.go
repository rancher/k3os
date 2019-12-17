package nodeupgrade

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	api "github.com/rancher/k3os/pkg/apis/k3os.io"
	apiv1 "github.com/rancher/k3os/pkg/apis/k3os.io/v1"
	"github.com/rancher/k3os/pkg/system"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func newUpgradeJob(spec apiv1.NodeUpgradeSpec, options Options) *batchv1.Job {
	var (
		backofLimit     = int32(2)
		deadlineSeconds = int64(600)
		privileged      = true
		upgradeKernel   = false
		upgradeRootFS   = true
	)
	debug, _ := strconv.ParseBool(os.Getenv("K3OS_DEBUG"))
	if ver, err := system.GetKernelVersion(); err != nil {
		logrus.Error(err)
	} else {
		upgradeKernel = ver.Current != ""
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels.Set{
				api.LabelUpgradeOperator: options.ControllerName,
				api.LabelUpgradeVersion:  spec.Version,
			},
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds: &deadlineSeconds,
			BackoffLimit:          &backofLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels.Set{
						api.LabelUpgradeOperator: options.ControllerName,
						api.LabelUpgradeVersion:  spec.Version,
					},
				},
				Spec: corev1.PodSpec{
					Affinity: upgradeJobAffinity(spec),
					Tolerations: []corev1.Toleration{{
						Key:      filepath.Join(corev1.LabelNamespaceSuffixNode, "unschedulable"),
						Operator: corev1.TolerationOpExists,
						Effect:   corev1.TaintEffectNoSchedule,
					}},
					RestartPolicy:      corev1.RestartPolicyNever,
					ServiceAccountName: options.ServiceAccountName,
					Volumes: []corev1.Volume{
						{Name: `bin`, VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{
							Path: "/bin", Type: hostPathType(corev1.HostPathDirectory),
						}}},
						{Name: `k3os-system`, VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{
							Path: system.RootPath(), Type: hostPathType(corev1.HostPathDirectory),
						}}},
						{Name: `lib`, VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{
							Path: "/lib", Type: hostPathType(corev1.HostPathDirectory),
						}}},
						{Name: `run-k3os`, VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{
							Path: system.StatePath(), Type: hostPathType(corev1.HostPathDirectory),
						}}},
						{Name: `sbin`, VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{
							Path: "/sbin", Type: hostPathType(corev1.HostPathDirectory),
						}}},
						{Name: `tmp`, VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{
							Path: "/tmp", Type: hostPathType(corev1.HostPathDirectory),
						}}},
						{Name: `usr`, VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{
							Path: "/usr", Type: hostPathType(corev1.HostPathDirectory),
						}}},
					},

					HostIPC: true,
					HostPID: true,

					InitContainers: []corev1.Container{{
						Image: "rancher/pause:3.1",
						Name:  "cordon",
						Command: []string{
							"kubectl", "cordon", spec.NodeName,
						},
						VolumeMounts: []corev1.VolumeMount{
							{Name: "k3os-system", MountPath: system.RootPath()},
							{Name: "tmp", MountPath: "/tmp"},
							{Name: "usr", MountPath: "/usr", ReadOnly: true},
						},
					}, {
						SecurityContext: &corev1.SecurityContext{
							Privileged: &privileged,
						},
						Image: fmt.Sprintf("rancher/k3os:%s", spec.Version),
						Name:  "upgrade",
						Command: []string{
							"k3os", "upgrade",
							"--source", "/k3os/system",
							"--destination", "/mnt/k3os/system",
							"--remount",
							"--kernel=" + strconv.FormatBool(upgradeKernel),
							"--rootfs=" + strconv.FormatBool(upgradeRootFS),
							"--sync",
						},
						Env: []corev1.EnvVar{
							{Name: "K3OS_DEBUG", Value: fmt.Sprintf("%v", debug)},
							//{Name: "PATH", Value: "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
						},
						VolumeMounts: []corev1.VolumeMount{
							{Name: "bin", MountPath: "/bin", ReadOnly: true},
							{Name: "k3os-system", MountPath: filepath.Join("/mnt", system.RootPath())},
							{Name: "lib", MountPath: "/lib", ReadOnly: true},
							{Name: "run-k3os", MountPath: system.StatePath()},
							{Name: "sbin", MountPath: "/sbin", ReadOnly: true},
							{Name: "tmp", MountPath: "/tmp"},
							{Name: "usr", MountPath: "/usr", ReadOnly: true},
						},
					}},

					Containers: []corev1.Container{{
						SecurityContext: &corev1.SecurityContext{
							Privileged: &privileged,
							Capabilities: &corev1.Capabilities{
								Add: []corev1.Capability{
									corev1.Capability("CAP_SYS_BOOT"),
								},
							},
						},
						Image: "rancher/pause:3.1",
						Name:  "reboot",
						// nsenter -m -u -i -n -p -t 1 -- reboot
						Command: []string{
							"nsenter", "-m", "-u", "-i", "-n", "-p", "-t", "1", "reboot",
						},
						VolumeMounts: []corev1.VolumeMount{
							{Name: "bin", MountPath: "/bin", ReadOnly: true},
							{Name: "k3os-system", MountPath: system.RootPath()},
							{Name: "lib", MountPath: "/lib", ReadOnly: true},
							{Name: "run-k3os", MountPath: system.StatePath()},
							{Name: "sbin", MountPath: "/sbin", ReadOnly: true},
							{Name: "tmp", MountPath: "/tmp"},
							{Name: "usr", MountPath: "/usr", ReadOnly: true},
						},
					}},
				},
			},
		},
	}
	if spec.Drain != nil {
		command := []string{"kubectl", "drain", spec.NodeName, "--pod-selector", `!` + api.LabelUpgradeOperator}
		if spec.Drain.IgnoreDaemonSets {
			command = append(command, "--delete-local-data")
		}
		if spec.Drain.DeleteLocalData {
			command = append(command, "--ignore-daemonsets")
		}
		if spec.Drain.Force {
			command = append(command, "--force")
		}
		if spec.Drain.Timeout != nil {
			command = append(command, "--timeout", spec.Drain.Timeout.String())
		}
		if spec.Drain.GracePeriod != nil {
			command = append(command, "--grace-period", strconv.FormatInt(int64(*spec.Drain.GracePeriod), 10))
		}
		job.Spec.Template.Spec.InitContainers = append(job.Spec.Template.Spec.InitContainers, corev1.Container{
			Image:        "rancher/pause:3.1",
			Name:         "drain",
			Command:      append(command),
			VolumeMounts: job.Spec.Template.Spec.InitContainers[0].VolumeMounts,
		})
	}
	return job
}

func upgradeJobAffinity(spec apiv1.NodeUpgradeSpec) *corev1.Affinity {
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{{
					MatchExpressions: []corev1.NodeSelectorRequirement{{
						Key:      corev1.LabelHostname,
						Operator: corev1.NodeSelectorOpIn,
						Values: []string{
							spec.NodeName,
						},
					}},
				}},
			},
		},
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{{
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{{
						Key:      api.LabelVersion,
						Operator: metav1.LabelSelectorOpExists,
					}},
				},
				TopologyKey: corev1.LabelHostname,
			}},
		},
	}
}

func hostPathType(hpt corev1.HostPathType) *corev1.HostPathType {
	return &hpt
}
