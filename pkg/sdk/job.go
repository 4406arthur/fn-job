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

//WebhookSetting ...
type WebhookSetting struct {
	Endpoint string `json:"endpoint"`
	Payload  string `json:"payload"`
}

//GenJobSpec used gen a job resource defined
func GenJobSpec(jobID, image string, cf *ConfigSetting, entryPoint []string, command []string, labels map[string]string, webhookSetting *WebhookSetting) *batchv1.Job {
	// TODO: retry should be changeable ?
	backoffLimit := int32(0)
	ttlSecondsAfterFinished := int32(3600)
	podSpec := apiv1.PodSpec{
		Containers: []apiv1.Container{
			{
				Name:    "smile-job",
				Image:   image,
				Command: entryPoint,
				Args:    command,
				// Resources: apiv1.ResourceRequirements{
				// 	Limits: apiv1.ResourceList{
				// 		"cpu":    resource.MustParse(cpuLimit),
				// 		"memory": resource.MustParse(memLimit),
				// 	},
				// 	Requests: apiv1.ResourceList{
				// 		"cpu":    resource.MustParse(cpuReq),
				// 		"memory": resource.MustParse(memReq),
				// 	},
				// },
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
	}

	//case for need injection configmap
	if cf != nil {
		//declare a volume for configmap
		podSpec.Volumes = []apiv1.Volume{
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
		}
		//volumn mount
		podSpec.Containers[0].VolumeMounts = []apiv1.VolumeMount{
			{
				Name:      "cofig-vol",
				MountPath: cf.ConfigDir,
			},
		}
		//pass a env var for app know configmap dir path
		podSpec.Containers[0].Env = append(podSpec.Containers[0].Env, apiv1.EnvVar{
			Name:  "CONFIGDIR",
			Value: cf.ConfigDir,
		})
	}

	objectMeta := metav1.ObjectMeta{
		Name: jobID,
	}

	if webhookSetting != nil {
		//injection into annotations
		objectMeta.Annotations = decodeWebhookConfig(webhookSetting)
		labels["webhook-enable"] = "true"
	} else {
		labels["webhook-enable"] = "false"
	}

	return &batchv1.Job{
		ObjectMeta: objectMeta,
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			//FEATURE STATE: Kubernetes v1.12 [alpha]
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: podSpec,
			},
		},
	}
}

func decodeWebhookConfig(ws *WebhookSetting) map[string]string {
	m := make(map[string]string)
	m["webhook-endpoint"] = ws.Endpoint
	m["webhook-payload"] = ws.Payload
	return m
}
