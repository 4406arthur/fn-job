package function

import (
	"context"
	"testing"

	"github.com/4406arthur/fn-job/pkg/sdk"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestGeneralCase(t *testing.T) {
	assert := assert.New(t)
	//testing data
	rq := Payload{
		JobID:      "64197512241",
		Image:      "busybox:latest",
		Namespace:  "default",
		EntryPoint: []string{"echo"},
		Command:    []string{"hello"},
	}
	err := validate.Struct(rq)
	assert.NoError(err, "checking nil error after validate rq")

	labels := map[string]string{"category": "mlaas-job"}
	jobSpec := sdk.GenJobSpec(rq.JobID, rq.Image, rq.Config, rq.EntryPoint, rq.Command, labels, rq.Webhook)
	assert.Equal(jobSpec.Spec.Template.Spec.Containers[0].Image, rq.Image)
	// assert.Equal(jobSpec.Spec.Template.Namespace, rq.Namesapce)
	assert.Equal(jobSpec.Spec.Template.Spec.Containers[0].Command, rq.EntryPoint)
	assert.Equal(jobSpec.Spec.Template.Spec.Containers[0].Args, rq.Command)

	kubeCli := testclient.NewSimpleClientset()
	_, err = kubeCli.BatchV1().Jobs(rq.Namespace).Create(
		context.TODO(),
		jobSpec,
		metav1.CreateOptions{},
	)
	assert.NoError(err, "checking nil error after kube-client apply")
}

func TestWebhookCase(t *testing.T) {
	assert := assert.New(t)
	//testing data
	rq := Payload{
		JobID:      "64197512241",
		Image:      "busybox:latest",
		Namespace:  "default",
		EntryPoint: []string{"echo"},
		Command:    []string{"hello"},
		Webhook: &sdk.WebhookSetting{
			Endpoint: "http://myapp.org/webhook/123461",
			Payload:  "{\"status:\": 200}",
		},
	}
	err := validate.Struct(rq)
	assert.NoError(err, "checking nil error after validate rq")
	labels := map[string]string{"category": "mlaas-job"}
	jobSpec := sdk.GenJobSpec(rq.JobID, rq.Image, rq.Config, rq.EntryPoint, rq.Command, labels, rq.Webhook)
	assert.Equal(jobSpec.ObjectMeta.Annotations["webhook-endpoint"], rq.Webhook.Endpoint)
	assert.Equal(jobSpec.ObjectMeta.Annotations["webhook-payload"], rq.Webhook.Payload)

	kubeCli := testclient.NewSimpleClientset()
	_, err = kubeCli.BatchV1().Jobs(rq.Namespace).Create(
		context.TODO(),
		jobSpec,
		metav1.CreateOptions{},
	)
	assert.NoError(err, "checking nil error after kube-client apply")
}

// Most resource types require a name that can be used as a DNS subdomain name as defined in RFC 1123. This means the name must
// contain no more than 253 characters
// contain only lowercase alphanumeric characters, '-' or '.'
// start with an alphanumeric character
// end with an alphanumeric character
func TestFormatErrorCase(t *testing.T) {
	assert := assert.New(t)
	//testing data
	rq := Payload{
		JobID:      "1213_41412",
		Image:      "busybox:latest",
		Namespace:  "default",
		EntryPoint: []string{"echo"},
		Command:    []string{"hello"},
	}

	err := validate.Struct(rq)
	assert.Error(err, "validate payload")
}
