package cmd

import (
	"flag"
	"fmt"

	"os"
)

const (
	NAMESPACE                      string = "devworkspace-controller" //TODO: Allow option to specify namespace
	usage                          string = "Takes as input an existing DevWorkspace and the path to a Devfile and prints to stdout a DevWorkspace object (which can be optionally applied to the cluster), identical to the orginal one, but with the template replaced by the Devfile content (with a few gotchas).\n\nUsage:\n  dw-update [options]\n\nOptions:\n  -d, --devfile=[]:\n    The file that contains the new devfile that is going to be applied.\n  -w, --devworkspace=[]:\n    The name of the original DevWorkspace object that is going to be used to create the new DevWorkspace.\n  -c, --cluster-mode=[true,false]\n    A boolean indicating whether the DevWorkspace on the cluster should be updated with the new DevWorkspace.\n  -f, --fetch=[true,false]\n    A boolean indicating whether the given DevWorkspace should be fetched by it's name on the cluster.\n"
	devFileArgHelpMessage          string = "The file that contains the new devfile that is going to be applied."
	devworkspaceHelpMessage        string = "The name of the original DevWorkspace object that is going to be used to create the new DevWorkspace. Requires --cluster-mode=true"
	updateClusterObjectHelpMessage string = "Whether the DevWorkspace object on the cluster should be updated with the new DevWorkspace. Requires --cluster-mode=true"
	fetchFromClusterHelpMessage    string = "Whether the given DevWorkspace should be fetched by it's name on the cluster."
)

type Options struct {
	DevfilePath         string
	DevWorkspaceName    string
	UpdateClusterObject bool
	FetchFromCluster    bool
}

func ParseArgs() *Options {
	devfilePath := flag.String("d", "", devFileArgHelpMessage)
	flag.StringVar(devfilePath, "devfile", *devfilePath, devFileArgHelpMessage)

	devworkspaceName := flag.String("w", "", devworkspaceHelpMessage)
	flag.StringVar(devworkspaceName, "devworkspace", *devworkspaceName, devworkspaceHelpMessage)

	updateClusterObject := flag.Bool("u", false, updateClusterObjectHelpMessage)
	flag.BoolVar(updateClusterObject, "update", *updateClusterObject, updateClusterObjectHelpMessage)

	fetchFromCluster := flag.Bool("f", true, fetchFromClusterHelpMessage)
	flag.BoolVar(fetchFromCluster, "fetch-from-cluster", *fetchFromCluster, fetchFromClusterHelpMessage)

	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}

	flag.Parse()

	if *devfilePath == "" {
		fmt.Println("A path to a devfile must be given.")
		os.Exit(1)
	}

	if *devworkspaceName == "" && !*fetchFromCluster {
		fmt.Println("Must provide a devworkspace name in order to fetch it from the cluster. Provide a devworkspace name with -w or --devworkspace")
		os.Exit(1)
	}

	if *devworkspaceName == "" {
		// Read devworkspace from stdin
		fmt.Println("The name of the devworkspace you want to update must be given.")
		os.Exit(1)
	}

	return &Options{
		DevfilePath:         *devfilePath,
		DevWorkspaceName:    *devworkspaceName,
		UpdateClusterObject: *updateClusterObject,
		FetchFromCluster:    *fetchFromCluster,
	}
}