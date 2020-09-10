package sdk

import (
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//ConfigSetting Support job mount a config from configmap
type ConfigSetting struct {
	ConfigMapRef string `json:"configMapRef"`
	ConfigDir    string `json:"configDir"`
}

//GenJobSpec used gen a job resource defined
func GenJobSpec(jobID, image string, entryPoint []string, command []string, labels map[string]string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobID,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "retrain-job",
							Image:   image,
							Command: entryPoint,
							Args:    command,
							Env: []apiv1.EnvVar{
								{
									Name: "MY_POD_NAMESPACE",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}
}

//GenAdvanceJobSpec used gen a job resource defined
func GenAdvanceJobSpec(jobID, image string, cf *ConfigSetting, entryPoint []string, command []string, labels map[string]string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobID,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "retrain-job",
							Image:   image,
							Command: entryPoint,
							Args:    command,
							Env: []apiv1.EnvVar{
								{
									Name: "MY_POD_NAMESPACE",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name:  "CONFIGDIR",
									Value: cf.ConfigDir,
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "cofig-vol",
									MountPath: cf.ConfigDir,
								},
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
					Volumes: []apiv1.Volume{
						{
							Name: "cofig-vol",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: cf.ConfigMapRef,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
