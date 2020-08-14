package function

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//Payload for http request
type Payload struct {
	Job        string   `json:"job"`
	Image      string   `json:"image"`
	EntryPoint []string `json:"entryPoint"`
	Command    []string `json:"command"`
}

//Handle ...
func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		// read request payload
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		input = reqBody
	}

	var rq Payload
	err := json.Unmarshal(input, &rq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}
	namespace := "mlaas-job"
	if s := os.Getenv("JOB_NAMESPACE"); len(s) > 0 {
		namespace = s
	}

	kubeCli, _ := newK8sCli("", "")
	_, err = kubeCli.BatchV1().Jobs(namespace).Create(
		context.TODO(),
		genJobSpec(rq.Job, rq.Image, rq.EntryPoint, rq.Command),
		metav1.CreateOptions{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

func newK8sCli(massterURL, kubeconfig string) (*kubernetes.Clientset, error) {

	if kubeconfig == "" && massterURL == "" {
		//klog.Warningf("Neither --kubeconfig nor --master was specified.  Using the inClusterConfig.  This might not work.")
		kubeconfig, err := restclient.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(kubeconfig)
		return clientset, nil
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags(massterURL, kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, nil
}

func genJobSpec(jobName, image string, entryPoint []string, command []string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    jobName,
							Image:   image,
							Command: entryPoint,
							Args:    command,
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}
}
