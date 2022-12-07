package main

import (
	"fmt"
	"path/filepath"

	"os"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	clusterClient "github.com/AObuchow/dw-update/customClient"
	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const NAMESPACE string = "devworkspace-controller"
const helpMessage string = "Takes as input an existing DevWorkspace and the path to a Devfile and prints to stdout a DevWorkspace object, identical to the orginal one, but with the template replaced by the Devfile content (with a few gotchas).\nUsage:\n  dw-update [options]\nOptions:  -d, --devfile=[]:\n    The file that contains the new devfile that is going to be applied.\n  -w, --devworksapce=[]:\n    The name of the original DevWorksapce object that is going to be used to create the new ."

func loadDevfileOrPanic(filePath string) dw.Devfile {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var devfile dw.Devfile
	if err := yaml.Unmarshal(bytes, &devfile); err != nil {
		panic(err)
	}
	return devfile
}

func main() {

	argsWithoutProg := os.Args[1:]

	devfilePath := ""
	devworkspaceName := ""

	// TODO: Improve commandline args handling, do error checking etc
	for i, arg := range argsWithoutProg {
		if arg == "-d" || arg == "--devfile" {
			devfilePath = argsWithoutProg[i+1]
		}

		if arg == "-w" || arg == "--devworkspace" {
			devworkspaceName = argsWithoutProg[i+1]
		}

		if arg == "-h" || arg == "--help" {
			fmt.Println(helpMessage)
			os.Exit(0)
		}
	}

	if devfilePath == "" {
		fmt.Println("A path to a devfile must be given.")
		os.Exit(1)
	}

	if devworkspaceName == "" {
		fmt.Println("The name of the devworkspace you want to update must be given.")
		os.Exit(1)
	}

	devfile := loadDevfileOrPanic(devfilePath)

	fmt.Println("Devfile is: ", devfilePath)
	fmt.Println("Devworkspace name is: ", devworkspaceName)

	fmt.Println("Devfile name: ", devfile.Metadata.Name)

	// TODO: Setup kube client depending on whether we're in a pod or running locally
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		panic(nil)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset

	client, err := clusterClient.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	dw, err := client.DevWorkspace(NAMESPACE).Get(devworkspaceName, metav1.GetOptions{})
	if err != nil && k8sErrors.IsNotFound(err) {
		panic(err)
	}

	// Get the devworkspace from cluster
	if dw != nil {
		fmt.Println("Found the dw with given name: " + dw.Name)

		// Preserve devworkspace spec.template.originalProjects
		originalProjects := dw.Spec.Template.Projects

		// take note of which spec.template.components have controller.devfile.io/merge-contribution: true attribute

		// todo use a map here to check for existence by component name when we it
		contributionNames := make(map[string]string)

		for _, component := range dw.Spec.Template.Components {
			if component.Attributes != nil {
				if component.Attributes.Exists("controller.devfile.io/merge-contribution") {
					if component.Attributes.GetBoolean("controller.devfile.io/merge-contribution", nil) {
						contributionNames[component.Name] = ""
					}
				}
			}
		}

		// Replace devworkspace spec.template with devfile content

		dw.Spec.Template = devfile.DevWorkspaceTemplateSpec

		// Retain original devworkspace projects
		dw.Spec.Template.Projects = originalProjects

		// for fun, append new projects..
		dw.Spec.Template.Projects = append(dw.Spec.Template.Projects, devfile.Projects...)

		// Retain merge contribution for components
		for _, component := range dw.Spec.Template.Components {
			if _, ok := contributionNames[component.Name]; ok {
				if !component.Attributes.Exists("controller.devfile.io/merge-contribution") {
					component.Attributes.PutBoolean("controller.devfile.io/merge-contribution", true)
				}

			}
		}

		// Update devworkspace on cluster
		_, err := client.DevWorkspace(NAMESPACE).Update(dw, metav1.UpdateOptions{})

		if err != nil {
			panic(err)
		}

	}

}
