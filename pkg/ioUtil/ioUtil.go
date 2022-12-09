package ioUtil // TODO: Find a better name?..

import (
	"fmt"

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

func LoadDevfileOrPanic(filePath string) dwv1alpha2.Devfile {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var devfile dwv1alpha2.Devfile
	if err := yaml.Unmarshal(bytes, &devfile); err != nil {
		panic(err)
	}
	return devfile
}
