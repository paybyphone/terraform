package command

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/cli"
)

func TestOutput(t *testing.T) {
	originalState := &terraform.State{
		Modules: []*terraform.ModuleState{
			&terraform.ModuleState{
				Path: []string{"root"},
				Outputs: map[string]*terraform.OutputState{
					"foo": &terraform.OutputState{
						Value: "bar",
						Type:  "string",
					},
				},
			},
		},
	}

	statePath := testStateFile(t, originalState)

	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{
		"-state", statePath,
		"foo",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}

	actual := strings.TrimSpace(ui.OutputWriter.String())
	if actual != "bar" {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestModuleOutput(t *testing.T) {
	originalState := &terraform.State{
		Modules: []*terraform.ModuleState{
			&terraform.ModuleState{
				Path: []string{"root"},
				Outputs: map[string]*terraform.OutputState{
					"foo": &terraform.OutputState{
						Value: "bar",
						Type:  "string",
					},
				},
			},
			&terraform.ModuleState{
				Path: []string{"root", "my_module"},
				Outputs: map[string]*terraform.OutputState{
					"blah": &terraform.OutputState{
						Value: "tastatur",
						Type:  "string",
					},
				},
			},
		},
	}

	statePath := testStateFile(t, originalState)

	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{
		"-state", statePath,
		"-module", "my_module",
		"blah",
	}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}

	actual := strings.TrimSpace(ui.OutputWriter.String())
	if actual != "tastatur" {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestModuleOutput_allOutputs(t *testing.T) {
	originalState := &terraform.State{
		Modules: []*terraform.ModuleState{
			&terraform.ModuleState{
				Path: []string{"root"},
				Outputs: map[string]*terraform.OutputState{
					"foo": &terraform.OutputState{
						Value: "bar",
					},
				},
			},
			&terraform.ModuleState{
				Path: []string{"root", "my_module"},
				Outputs: map[string]*terraform.OutputState{
					"blah": &terraform.OutputState{
						Value: "tastatur",
					},
					"baz": &terraform.OutputState{
						Value: "qux",
					},
				},
			},
		},
	}

	statePath := testStateFile(t, originalState)

	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{
		"-state", statePath,
		"-module", "my_module",
		"blah",
	}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}

	expectedOutput := strings.Split("\n", "blah = tastatur\nbaz = qux\n")
	sort.Strings(expectedOutput)

	output := strings.Split("\n", ui.OutputWriter.String())
	sort.Strings(output)
	if reflect.DeepEqual(output, expectedOutput) != true {
		t.Fatalf("Expected output: %#v\ngiven: %#v", expectedOutput, output)
	}
}

func TestMissingModuleOutput(t *testing.T) {
	originalState := &terraform.State{
		Modules: []*terraform.ModuleState{
			&terraform.ModuleState{
				Path: []string{"root"},
				Outputs: map[string]*terraform.OutputState{
					"foo": &terraform.OutputState{
						Value: "bar",
						Type:  "string",
					},
				},
			},
		},
	}

	statePath := testStateFile(t, originalState)

	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{
		"-state", statePath,
		"-module", "not_existing_module",
		"blah",
	}

	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
}

func TestOutput_badVar(t *testing.T) {
	originalState := &terraform.State{
		Modules: []*terraform.ModuleState{
			&terraform.ModuleState{
				Path: []string{"root"},
				Outputs: map[string]*terraform.OutputState{
					"foo": &terraform.OutputState{
						Value: "bar",
						Type:  "string",
					},
				},
			},
		},
	}

	statePath := testStateFile(t, originalState)

	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{
		"-state", statePath,
		"bar",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
}

func TestOutput_blank(t *testing.T) {
	originalState := &terraform.State{
		Modules: []*terraform.ModuleState{
			&terraform.ModuleState{
				Path: []string{"root"},
				Outputs: map[string]*terraform.OutputState{
					"foo": &terraform.OutputState{
						Value: "bar",
						Type:  "string",
					},
					"name": &terraform.OutputState{
						Value: "john-doe",
						Type:  "string",
					},
				},
			},
		},
	}

	statePath := testStateFile(t, originalState)

	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{
		"-state", statePath,
		"",
	}

	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}

	expectedOutput := "foo = bar\nname = john-doe\n"
	output := ui.OutputWriter.String()
	if output != expectedOutput {
		t.Fatalf("Expected output: %#v\ngiven: %#v", expectedOutput, output)
	}
}

func TestOutput_manyArgs(t *testing.T) {
	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{
		"bad",
		"bad",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: \n%s", ui.OutputWriter.String())
	}
}

func TestOutput_noArgs(t *testing.T) {
	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: \n%s", ui.OutputWriter.String())
	}
}

func TestOutput_noState(t *testing.T) {
	originalState := &terraform.State{}
	statePath := testStateFile(t, originalState)

	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{
		"-state", statePath,
		"foo",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
}

func TestOutput_noVars(t *testing.T) {
	originalState := &terraform.State{
		Modules: []*terraform.ModuleState{
			&terraform.ModuleState{
				Path:    []string{"root"},
				Outputs: map[string]*terraform.OutputState{},
			},
		},
	}

	statePath := testStateFile(t, originalState)

	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{
		"-state", statePath,
		"bar",
	}
	if code := c.Run(args); code != 1 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}
}

func TestOutput_stateDefault(t *testing.T) {
	originalState := &terraform.State{
		Modules: []*terraform.ModuleState{
			&terraform.ModuleState{
				Path: []string{"root"},
				Outputs: map[string]*terraform.OutputState{
					"foo": &terraform.OutputState{
						Value: "bar",
						Type:  "string",
					},
				},
			},
		},
	}

	// Write the state file in a temporary directory with the
	// default filename.
	td, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	statePath := filepath.Join(td, DefaultStateFilename)

	f, err := os.Create(statePath)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = terraform.WriteState(originalState, f)
	f.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Change to that directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if err := os.Chdir(filepath.Dir(statePath)); err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Chdir(cwd)

	ui := new(cli.MockUi)
	c := &OutputCommand{
		Meta: Meta{
			ContextOpts: testCtxConfig(testProvider()),
			Ui:          ui,
		},
	}

	args := []string{
		"foo",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: \n%s", ui.ErrorWriter.String())
	}

	actual := strings.TrimSpace(ui.OutputWriter.String())
	if actual != "bar" {
		t.Fatalf("bad: %#v", actual)
	}
}
