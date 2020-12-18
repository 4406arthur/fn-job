package function

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/4406arthur/fn-job/pkg/sdk"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Payload for http request
type Payload struct {
	JobID      string              `json:"jobID"`
	Image      string              `json:"image"`
	Config     *sdk.ConfigSetting  `json:"config,omitempty"`
	Namesapce  string              `json:"namespace"`
	EntryPoint []string            `json:"entryPoint"`
	Command    []string            `json:"command"`
	Webhook    *sdk.WebhookSetting `json:"webhook,omitempty"`
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

	kubeCli, _ := sdk.NewK8sCli("", "")
	//give a specific label
	labels := map[string]string{"category": "mlaas-job"}

	var jobSpec *v1.Job
	jobSpec = sdk.GenJobSpec(rq.JobID, rq.Image, rq.Config, rq.EntryPoint, rq.Command, labels, rq.Webhook)

	_, err = kubeCli.BatchV1().Jobs(rq.Namesapce).Create(
		context.TODO(),
		jobSpec,
		metav1.CreateOptions{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}
