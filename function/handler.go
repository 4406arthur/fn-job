package function

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/4406arthur/fn-job/pkg/sdk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Payload for http request
type Payload struct {
	JobID      string   `json:"jobID"`
	Image      string   `json:"image"`
	Namesapce  string   `json:"namespace"`
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

	kubeCli, _ := sdk.NewK8sCli("", "")
	//give a specific label
	labels := map[string]string{"category": "mlaas-job"}
	_, err = kubeCli.BatchV1().Jobs(rq.Namesapce).Create(
		context.TODO(),
		sdk.GenJobSpec(rq.JobID, rq.Image, rq.EntryPoint, rq.Command, labels),
		metav1.CreateOptions{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}
