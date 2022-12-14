package customClient

import (
	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ExampleV1Alpha1Interface interface {
	Projects(namespace string) DevworkspaceClientInterface
}

type ExampleV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ExampleV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: dw.SchemeGroupVersion.Group, Version: dw.SchemeGroupVersion.Version}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	dw.AddToScheme(scheme.Scheme)

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ExampleV1Alpha1Client{restClient: client}, nil
}

func (c *ExampleV1Alpha1Client) DevWorkspace(namespace string) DevworkspaceClientInterface {
	return &devworkspaceClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
