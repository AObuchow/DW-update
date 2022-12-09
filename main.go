package main

import (
	"fmt"

	"os"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"

	cmd "github.com/AObuchow/dw-update/pkg/cmd"
	clusterClient "github.com/AObuchow/dw-update/pkg/customClient"
	io "github.com/AObuchow/dw-update/pkg/ioUtil"
	kube "github.com/AObuchow/dw-update/pkg/kube"
	dwv1alpha2 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NAMESPACE string = "devworkspace-controller" //TODO: Allow option to specify namespace
)

func main() {
	opts := cmd.ParseArgs()
	devfile := io.LoadDevfileOrPanic(opts.DevfilePath)

	config, err := kube.GetKubeConfig()
	if err != nil {
		// TODO: Use a logger
		panic(err.Error())
	}

	// TODO: Remove these, for debug purposes
	fmt.Println("Devfile is: ", opts.DevfilePath)
	fmt.Println("Devworkspace name is: ", opts.DevWorkspaceName)
	fmt.Println("Devfile name: ", devfile.Metadata.Name)

	// create the clientset
	client, err := clusterClient.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	dw, err := client.DevWorkspace(NAMESPACE).Get(opts.DevWorkspaceName, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			fmt.Fprintf(os.Stderr, "Couldn't find DevWorkspace with name %s on the cluster", opts.DevWorkspaceName)
			os.Exit(1)
		}
		panic(err)
	}

	// Get the devworkspace from cluster
	if dw != nil {
		fmt.Println("Found the dw with given name: " + dw.Name)
		dw = updateDevWorkspace(*dw, devfile)

		if opts.UpdateClusterObject {
			// Update devworkspace on cluster
			_, err := client.DevWorkspace(NAMESPACE).Update(dw, metav1.UpdateOptions{})
			if err != nil {
				panic(err)
			}
			fmt.Printf("Succesfully updated devworkspace %s\n", dw.Name)
		}

		io.PrintDevWorkspace(dw)
	}
}

func updateDevWorkspace(dw dwv1alpha2.DevWorkspace, devfile dwv1alpha2.Devfile) (updatedDevWorkspace *dwv1alpha2.DevWorkspace) {
	updatedDevWorkspace = dw.DeepCopy()
	// Preserve original devworkspace spec.template.projects
	originalProjects := updatedDevWorkspace.Spec.Template.Projects

	// Find component with the controller.devfile.io/merge-contribution: true attribute
	mergeContributionComponent := ""
	for _, component := range updatedDevWorkspace.Spec.Template.Components {
		if component.Attributes != nil {
			if component.Attributes.Exists("controller.devfile.io/merge-contribution") {
				if component.Attributes.GetBoolean("controller.devfile.io/merge-contribution", nil) {
					mergeContributionComponent = component.Name
					break // There is only supposed to be one merge contribution component so we stop once we find it
				}
			}
		}
	}

	// Replace devworkspace spec.template with devfile content
	updatedDevWorkspace.Spec.Template = devfile.DevWorkspaceTemplateSpec

	// Retain original devworkspace projects
	// TODO: Append here so that the user can add more projects when updating devworkspace?
	updatedDevWorkspace.Spec.Template.Projects = originalProjects

	// Retain merge contribution attribute
	for _, component := range updatedDevWorkspace.Spec.Template.Components {
		if component.Name == mergeContributionComponent {
			if !component.Attributes.Exists("controller.devfile.io/merge-contribution") {
				component.Attributes.PutBoolean("controller.devfile.io/merge-contribution", true)
			}
		}
	}
	return updatedDevWorkspace
}
