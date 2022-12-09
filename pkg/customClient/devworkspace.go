package customClient

import (
	"context"

	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type DevworkspaceClientInterface interface {
	List(opts metav1.ListOptions) (*dw.DevWorkspaceList, error)
	Get(name string, opts metav1.GetOptions) (*dw.DevWorkspace, error)
	Create(*dw.DevWorkspace) (*dw.DevWorkspace, error)
	Update(*dw.DevWorkspace, metav1.UpdateOptions) (*dw.DevWorkspace, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type devworkspaceClient struct {
	restClient rest.Interface
	ns         string
}

const resource string = "devWorkspaces"

func (c *devworkspaceClient) List(opts metav1.ListOptions) (*dw.DevWorkspaceList, error) {
	result := dw.DevWorkspaceList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *devworkspaceClient) Get(name string, opts metav1.GetOptions) (*dw.DevWorkspace, error) {
	result := dw.DevWorkspace{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resource).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *devworkspaceClient) Update(devworkspace *dw.DevWorkspace, opts metav1.UpdateOptions) (*dw.DevWorkspace, error) {
	result := dw.DevWorkspace{}
	err := c.restClient.
		Put().
		Namespace(c.ns).
		Resource(resource).
		Body(devworkspace).
		Name(devworkspace.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *devworkspaceClient) Create(devworkspace *dw.DevWorkspace) (*dw.DevWorkspace, error) {
	result := dw.DevWorkspace{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(resource).
		Body(devworkspace).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *devworkspaceClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
