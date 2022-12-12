package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/AObuchow/dw-update/pkg/ioUtil"
	dwv1alpha2 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

const (
	usage string = `Takes as input an existing DevWorkspace (by it's name on the cluster, a YAML file or through stdin) and the path to a Devfile and prints to stdout a DevWorkspace object (which can be optionally applied to the cluster), identical to the orginal one, but with the template replaced by the Devfile content (with a few gotchas).

Usage:
	$ dw-update [options]

Example usages:
	$ dw-update -d <devfile> -w <devworkspace> # Takes in a Devfile and DevWorkspace path, outputs updated DevWorkspace
	$ dw-update -d <devfile> -n <devworkspace name> -f true # Takes in a Devfile path and a DevWorkspace name which is used to fetch a DevWorkspace from the cluster, outputs updated DevWorkspace.
	$ dw-update -d <devfile> # Takes in a Devfile path and a DevWorkspace from stdin, outputs updated DevWorkspace.
	$ dw-update -d <devfile> -w <devworkspace> -c true # Takes in a Devfile and DevWorkspace path, outputs updated DevWorkspace and updates the DevWorkspace on the cluster.
	
Options:
	-d, --devfile=[]:
	The path to the file that contains the new devfile that is going to be applied.
	-w, --devworkspace=[]:
	The path to the file that contains the original devworkspace.
	-n, --devworkspace-name=[]:
	The name of the original DevWorkspace object (on the cluster) that is going to be used to create the new DevWorkspace.
	Requires --fetch=true
	-c, --cluster-mode=[true,false]
	A boolean indicating whether the DevWorkspace on the cluster should be updated with the new DevWorkspace.
	-f, --fetch=[true,false]
	A boolean indicating whether the given DevWorkspace should be fetched by it's name on the cluster.
`
	devfileHelpMessage             string = "The path to the file that contains the new devfile that is going to be applied."
	devworkspacePathHelpMessage    string = "The path to the file that contains the original Devworkspace that is going to be used to create the new DevWorkspace."
	devworkspaceNameHelpMessage    string = "The name of the original DevWorkspace object (on the cluster) that is going to be used to create the new DevWorkspace. Requires --fetch=true."
	updateClusterObjectHelpMessage string = "Whether the DevWorkspace object on the cluster should be updated with the new DevWorkspace."
	fetchFromClusterHelpMessage    string = "Whether the given DevWorkspace should be fetched by it's name on the cluster."
)

type Options struct {
	DevfilePath         string
	DevWorkspaceName    string
	UpdateClusterObject bool
	FetchFromCluster    bool
	ParsedDevWorkspace  *dwv1alpha2.DevWorkspace
}

func ParseArgs() *Options {
	// TODO: No need for these extra variables, just create an Options struct right away?
	var parsedDW *dwv1alpha2.DevWorkspace = nil

	devfilePath := flag.String("d", "", devfileHelpMessage)
	flag.StringVar(devfilePath, "devfile", *devfilePath, devfileHelpMessage)

	devworkspacePath := flag.String("w", "", devfileHelpMessage)
	flag.StringVar(devworkspacePath, "devworkspace", *devworkspacePath, devfileHelpMessage)

	devworkspaceName := flag.String("n", "", devworkspacePathHelpMessage)
	flag.StringVar(devworkspaceName, "devworkspace-name", *devworkspaceName, devworkspacePathHelpMessage)

	updateClusterObject := flag.Bool("u", false, updateClusterObjectHelpMessage)
	flag.BoolVar(updateClusterObject, "update", *updateClusterObject, updateClusterObjectHelpMessage)
	// TODO: Get rid of this option. If you provide the devworkspace name, fetch from cluster should be enabled.
	fetchFromCluster := flag.Bool("f", false, fetchFromClusterHelpMessage)
	flag.BoolVar(fetchFromCluster, "fetch", *fetchFromCluster, fetchFromClusterHelpMessage)

	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}

	flag.Parse()

	if *devfilePath == "" {
		fmt.Fprintf(os.Stderr, "A path to a devfile must be given.\n")
		os.Exit(1)
	}

	if *devworkspaceName != "" && *devworkspacePath != "" {
		fmt.Fprintf(os.Stderr, "Cannot provide multiple DevWorkspace input methods. Specify either -w OR -n. Alternatively, omit both -w and -n and provide the DevWorkspace via stdin.\n")
		os.Exit(1)
	}

	if *devworkspacePath != "" {
		devworkspace, err := ioUtil.LoadDevWorkspace(*devworkspacePath)
		if err != nil {
			panic(err)
		}
		// TODO: This is very hacky and should be cleaned up. The devworkspace should be loaded in elsewhere. This function should deal exclusively with parsing and validating arguments.
		parsedDW = devworkspace
	}

	// TODO: Should stdin mode require an explicit flag argument?
	if *devworkspaceName == "" && *devworkspacePath == "" && !*fetchFromCluster {
		// Read devworkspace from stdin
		devworkspace, err := ioUtil.ParseDevWorkspaceStdin()
		if err != nil {
			panic(err)
		}
		parsedDW = devworkspace
	}

	if parsedDW == nil && *devworkspaceName == "" && *fetchFromCluster {
		fmt.Println("Must provide a devworkspace name in order to fetch it from the cluster. Provide a devworkspace name with -n or --devworkspace-name")
		os.Exit(1)
	}

	return &Options{
		DevfilePath:         *devfilePath,
		DevWorkspaceName:    *devworkspaceName,
		UpdateClusterObject: *updateClusterObject,
		FetchFromCluster:    *fetchFromCluster,
		ParsedDevWorkspace:  parsedDW,
	}
}
