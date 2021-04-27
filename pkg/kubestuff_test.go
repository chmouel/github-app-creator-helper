package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-github/v35/github"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNew(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fail()
	}

	r, err := os.Open("/tmp/a.json")
	if err != nil {
		t.Fail()
	}
	bits, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fail()
	}
	var manifest github.AppConfig
	json.Unmarshal(bits, &manifest)

	data := map[string][]byte{
		"application_id": []byte(fmt.Sprintf("%d", manifest.GetID())),
		"private.key":    []byte(manifest.GetPEM()),
		"webhook.secret": []byte(manifest.GetWebhookSecret()),
	}

	var secretName = "github-app-secret"

	if sec, _ := c.kubeClient.CoreV1().Secrets(c.Namespace).Get(secretName, metav1.GetOptions{}); sec != nil {
		fmt.Printf("We already have a secret called %s\n", secretName)
		t.Fail()
	}
	c.kubeClient.CoreV1().Secrets(c.Namespace).Create(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Data: data,
	})

}
