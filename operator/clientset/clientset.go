package clientset

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	cr "github.com/alexrodfe/golang-api/operator/deployment-objects"
)

type CustomClientInterface interface {
	DeploymentObejcts(namespace string) DeploymentObjectsInterface
}

type CustomClient struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*CustomClient, error) {
	// adds custom resources to scheme
	err := cr.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}

	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: cr.GroupName, Version: cr.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &CustomClient{restClient: client}, nil
}

func (c *CustomClient) DeploymentObejcts(namespace string) DeploymentObjectsInterface {
	return &deploymentObjectsClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
