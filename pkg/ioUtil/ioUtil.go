package ioUtil // TODO: Find a better name?..

import (
	"fmt"
	"io"
	"os"

	dwv1alpha2 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/yaml"
)

var yamlPrinter printers.YAMLPrinter = printers.YAMLPrinter{}

func PrintDevWorkspace(dw *dwv1alpha2.DevWorkspace) {
	fmt.Printf("Resulting DevWorkspace:\n\n\n")
	dw.GetObjectKind().SetGroupVersionKind(dwv1alpha2.SchemeGroupVersion.WithKind("DevWorkspace"))
	yamlPrinter.PrintObj(dw, os.Stdout)
}

func LoadDevfile(filePath string) (*dwv1alpha2.Devfile, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	devfile := &dwv1alpha2.Devfile{}
	if err := yaml.Unmarshal(bytes, &devfile); err != nil {
		return nil, err
	}
	return devfile, nil
}

func ParseDevWorkspaceStdin() (*dwv1alpha2.DevWorkspace, error) {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	devworkspace := &dwv1alpha2.DevWorkspace{}
	if err := yaml.Unmarshal(bytes, &devworkspace); err != nil {
		return nil, err
	}
	devworkspace.Name = devworkspace.ObjectMeta.Name

	PrintDevWorkspace(devworkspace)
	return devworkspace, nil
}
