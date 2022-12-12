package update

import (
	"os"
	"path/filepath"
	"testing"

	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

type testCase struct {
	Name   string     `json:"name,omitempty"`
	Input  testInput  `json:"input,omitempty"`
	Output testOutput `json:"output,omitempty"`
}

type testInput struct {
	Devfile      *dw.Devfile      `json:"devfile,omitempty"`
	DevWorkspace *dw.DevWorkspace `json:"devworkspace,omitempty"`
}

type testOutput struct {
	Devworkspace *dw.DevWorkspace `json:"devworkspace,omitempty"`
}

func loadTestCaseOrPanic(t *testing.T, testFilepath string) testCase {
	bytes, err := os.ReadFile(testFilepath)
	if err != nil {
		t.Fatal(err)
	}
	var test testCase
	if err := yaml.Unmarshal(bytes, &test); err != nil {
		t.Fatal(err)
	}
	return test
}

func loadAllTestCasesOrPanic(t *testing.T, fromDir string) []testCase {
	files, err := os.ReadDir(fromDir)
	if err != nil {
		t.Fatal(err)
	}
	var tests []testCase
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		tests = append(tests, loadTestCaseOrPanic(t, filepath.Join(fromDir, file.Name())))
	}
	return tests
}

func TestUpdateDevWorkspace(t *testing.T) {
	tests := loadAllTestCasesOrPanic(t, "testdata/")

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// sanity check that file is read correctly.
			assert.NotNil(t, tt.Input.Devfile, "Input does not define a devfile")
			assert.NotNil(t, tt.Input.DevWorkspace, "Input does not define a devworkspace") // TODO: Remove this when functionality is added to allow only taking a devfile input
			devfile := *tt.Input.Devfile
			devworkspace := *tt.Input.DevWorkspace

			updatedDevWorkspace := UpdateDevWorkspace(devworkspace, &devfile)
			assert.Equal(t, tt.Output.Devworkspace, updatedDevWorkspace, "Updated devworkspace does not match expected output devworkspace.")
		})
	}
}
