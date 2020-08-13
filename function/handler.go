package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/4406arthur/fn-job/pkg/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	namespace := os.Getenv("JOB_NAMESPACE")
	fmt.Printf("job deploy namespace is: %s", namespace)

	kubeCli, _ := apis.NewK8sCli("", "")
	_, err = kubeCli.BatchV1().Jobs("default").Create(
		context.TODO(),
		apis.GenJobSpec(rq.Job, rq.Image, namespace, rq.EntryPoint, rq.Command),
		metav1.CreateOptions{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}
