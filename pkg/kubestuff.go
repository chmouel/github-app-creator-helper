package pkg

import (
	"errors"

	routev1 "github.com/openshift/api/route/v1"
	routeclientset "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	kubeClient  k8s.Interface
	routeClient routeclientset.RouteV1Interface
	Namespace   string
}

func (c *Client) GetFirstTLSRoute() (*routev1.Route, error) {
	routeList, err := c.routeClient.Routes(c.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, value := range routeList.Items {
		if value.Spec.TLS != nil {
			return &value, nil
		}
	}

	return nil, errors.New("No TLS Route found.")
}

func NewClient() (*Client, error) {
	var client Client
	config, ns, err := kubeClient()
	if err != nil {
		return nil, err
	}
	client.Namespace = ns

	k8scs, err := k8s.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client.kubeClient = k8scs

	client.routeClient, err = routeclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func kubeClient() (*rest.Config, string, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		return nil, "", errors.New("Couldn't get kubeConfiguration namespace")
	}

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, "", errors.New("Parsing kubeconfig failed")
	}

	return config, namespace, nil
}
