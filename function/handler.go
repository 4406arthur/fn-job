package function

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/4406arthur/fn-job/pkg/sdk"
	"github.com/go-playground/validator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Payload for http request
type Payload struct {
	JobID      string              `json:"jobID" validate:"required,hostname_rfc1123"`
	Image      string              `json:"image" validate:"required"`
	Config     *sdk.ConfigSetting  `json:"config,omitempty"`
	Namespace  string              `json:"namespace" validate:"required,hostname_rfc1123"`
	EntryPoint []string            `json:"entryPoint" validate:"required"`
	Command    []string            `json:"command" validate:"required"`
	Webhook    *sdk.WebhookSetting `json:"webhook,omitempty"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

//Handle ...
func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		// read request payload
		reqBody, err := ioutil.ReadAll(r.Body)

		// log caller info
		log.Printf("caller info: %s %s %s\n", r.UserAgent(), r.RemoteAddr, reqBody)

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

	err = validate.Struct(rq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	kubeCli, _ := sdk.NewK8sCli("", "")
	//give a specific label
	labels := map[string]string{"category": "mlaas-job"}

	jobSpec := sdk.GenJobSpec(rq.JobID, rq.Image, rq.Config, rq.EntryPoint, rq.Command, labels, rq.Webhook)

	_, err = kubeCli.BatchV1().Jobs(rq.Namespace).Create(
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
