package clientset

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	cr "github.com/alexrodfe/golang-api/operator/deployment-objects"
)

type DeploymentObjectsInterface interface {
	List(opts metav1.ListOptions) (*cr.DeploymentObjectList, error)
	Get(name string, options metav1.GetOptions) (*cr.DeploymentObject, error)
	Create(*cr.DeploymentObject) (*cr.DeploymentObject, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type deploymentObjectsClient struct {
	restClient rest.Interface
	ns         string
}

func (c *deploymentObjectsClient) List(opts metav1.ListOptions) (*cr.DeploymentObjectList, error) {
	result := cr.DeploymentObjectList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("deployment-objects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *deploymentObjectsClient) Get(name string, opts metav1.GetOptions) (*cr.DeploymentObject, error) {
	result := cr.DeploymentObject{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("deployment-objects").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *deploymentObjectsClient) Create(project *cr.DeploymentObject) (*cr.DeploymentObject, error) {
	result := cr.DeploymentObject{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("deployment-objects").
		Body(project).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *deploymentObjectsClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("deployment-objects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
