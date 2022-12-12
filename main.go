package main

import (
	"fmt"

	"os"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"

	cmd "github.com/AObuchow/dw-update/pkg/cmd"
	clusterClient "github.com/AObuchow/dw-update/pkg/customClient"
	io "github.com/AObuchow/dw-update/pkg/ioUtil"
	kube "github.com/AObuchow/dw-update/pkg/kube"
	update "github.com/AObuchow/dw-update/pkg/update"
	dwv1alpha2 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NAMESPACE string = "devworkspace-controller" //TODO: Allow option to specify namespace
)

func main() {
	opts := cmd.ParseArgs()
	devfile, err := io.LoadDevfile(opts.DevfilePath)
	if err != nil {
		panic(err)
	}

	var client *clusterClient.ExampleV1Alpha1Client = nil
	if opts.UpdateClusterObject || opts.FetchFromCluster {
		config, err := kube.GetKubeConfig()
		if err != nil {
			// TODO: Use a logger
			panic(err.Error())
		}

		client, err = clusterClient.NewForConfig(config)
		if err != nil {
			panic(err)
		}
	}

	// TODO: Remove these, for debug purposes
	fmt.Println("Devfile is: ", opts.DevfilePath)
	fmt.Println("Devworkspace name is: ", opts.DevWorkspaceName)
	fmt.Println("Devfile name: ", devfile.Metadata.Name)

	var dw *dwv1alpha2.DevWorkspace = nil
	if opts.FetchFromCluster {
		dw, err = client.DevWorkspace(NAMESPACE).Get(opts.DevWorkspaceName, metav1.GetOptions{})
		if err != nil {
			if k8sErrors.IsNotFound(err) {
				fmt.Fprintf(os.Stderr, "Couldn't find DevWorkspace with name %s on the cluster", opts.DevWorkspaceName)
				os.Exit(1)
			}
			panic(err)
		}
	} else {
		dw = opts.ParsedDevWorkspace
	}

	// Get the devworkspace from cluster
	if dw != nil {
		fmt.Println("Found the dw with given name: " + dw.Name)
		dw = update.UpdateDevWorkspace(*dw, devfile)

		if opts.UpdateClusterObject {
			// Update devworkspace on cluster
			_, err := client.DevWorkspace(NAMESPACE).Update(dw, metav1.UpdateOptions{})
			if err != nil {
				panic(err)
			}
			fmt.Printf("Succesfully updated devworkspace %s\n", dw.Name)
		}

		io.PrintDevWorkspace(dw)
	} else {
		fmt.Fprint(os.Stderr, "No devworkspace provided.", opts.DevWorkspaceName)
		os.Exit(1)
	}
}
