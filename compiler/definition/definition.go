package definition

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"gomate.io/gomate/logging"
)

var emptyLineRexexp = regexp.MustCompile("^[\t ]*$")

// Definition represents a parsed step definition, typically located in step_definition folder underneath features folder.
type Definition struct {
	imports []string
	funcs   []string
}

type stepDefinitions []Definition

// Definitions has been composed by multiple Definion instances, but does also contain
// logic to decide what and where to execute, plus rules how to clean up binary result.
type Definitions struct {
	defs    stepDefinitions
	tmpDir  string
	command string
	removed bool
}

func (definitions stepDefinitions) Code() string {
	imports := []string{}
	funcs := []string{}

	for _, definition := range definitions {
		imports = append(imports, definition.imports...)
		funcs = append(funcs, definition.funcs...)
	}

	return fmt.Sprintf(snippet, strings.Join(imports, "\n"), strings.Join(funcs, "\n"))
}

// Code method generates source code based on a step definition.
// The step definition are located into a go function named setup, this function sets
// up all definition just before the parsement of the feature file that has been supplied as first
// argument on commandline. Second argument must be text string "true" or "false", that enables
// and disables pretty print i.e., print to STDOUT with or without color.
func (definition Definition) Code() string {
	return fmt.Sprintf(
		snippet,
		strings.Join(definition.imports, "\n"),
		strings.Join(definition.funcs, "\n"),
	)
}

// Code method generates source code based on step definitions.
// The composited step definitions are located into a go function named setup, this function sets
// up all definition just before the parsement of the feature file that has been supplied as first
// argument on commandline. Second argument must be text string "true" or "false", that enables
// and disables pretty print i.e., print to STDOUT with or without color.
func (definitions Definitions) Code() string {
	return definitions.defs.Code()
}

// Run takes a DSL written feature from io.Reader and supply the data into precompiled behaviour code.
// After execution of the binary, the result are written to STDOUT. Method respects global.PPrint to
// enable/disable pretty print i.e., colors enabled.
func (definitions Definitions) Run(features io.Reader, pprint bool) {
	if definitions.removed {
		logging.Info("Compiled behaviour binary file has been removed")
		return
	}

	featureLines, err := ioutil.ReadAll(features)

	if err != nil {
		logging.Fatal(err.Error())
	}

	gorun := exec.Command(definitions.command, strconv.FormatBool(pprint)) // #nosec
	if stdin, err := gorun.StdinPipe(); err != nil {
		logging.Fatal(err.Error())
	} else if n, err := stdin.Write(featureLines); err != nil {
		logging.Fatal(err.Error())
	} else if n != len(featureLines) {
		logging.Fatal("Behaviour binary file was not able to read all defined features")
	} else if err := stdin.Close(); err != nil {
		logging.Fatal(err.Error())
	}

	if output, err := gorun.CombinedOutput(); err != nil {
		logging.Fatal(err.Error())
	} else {
		logging.Info(string(output))
	}
}

// NewDefinition reads and parse "in" assuming that it contains content from a step definition file.
// Lines defining package names are omitted from resulting Definition instance.
func NewDefinition(in io.Reader) Definition {
	code, err := ioutil.ReadAll(in)
	if err != nil {
		logging.Fatal(err.Error())
	}

	imports := []string{}
	funcs := []string{}

	for _, row := range strings.Split(string(code), "\n") {
		if strings.HasPrefix(row, "import ") {
			imports = append(imports, row)
		} else if strings.HasPrefix(row, "package ") {
			continue // Package removed, "main" added later
		} else if emptyLineRexexp.MatchString(row) {
			continue // Empty lines removed
		} else {
			funcs = append(funcs, row)
		}
	}

	return Definition{imports, funcs}
}

// NewDefinitions reads and parse "defs" assuming that each element contains content from a step definition file.
// Lines defining package names are omitted from resulting Definition instance.
//
// NOTE: It's callers responsibility to call Remove method before Definitions instance are garbage
// collected.
func NewDefinitions(defs []io.Reader, forensic bool) Definitions {
	definitions := stepDefinitions{}

	for _, reader := range defs {
		definitions = append(definitions, NewDefinition(reader))
	}

	dir, cmd := definitions.compile(forensic)
	return Definitions{definitions, dir, cmd, false}
}

// Remove will remove temporary file containing the generated step definition
// After this method has been called, it's no longer possible to execute Run
// method.
func (definitions Definitions) Remove() {
	definitions.removed = true // Don't allow Run-calls from now on

	if err := os.RemoveAll(definitions.tmpDir); err != nil {
		logging.Err(err.Error())
		return
	}
}

func (definitions stepDefinitions) compile(forensic bool) (string, string) {
	var err error
	var output []byte

	dir, testCode, testFile := definitions.store(forensic)

	goimport := exec.Command("goimports", "-w=true", testCode)       // #nosec
	gofmt := exec.Command("go", "fmt", testCode)                     // #nosec
	gobuild := exec.Command("go", "build", "-o", testFile, testCode) // #nosec

	if err = goimport.Run(); err != nil {
		logging.Err(err.Error())
	}

	if err = gofmt.Run(); err != nil {
		logging.Err(err.Error())
	}

	if output, err = gobuild.CombinedOutput(); err != nil {
		logging.Err(string(output))
	}

	return dir, testFile
}

func (definitions stepDefinitions) store(forensic bool) (dir, testCode, testFile string) {
	var err error

	if dir, err = ioutil.TempDir("", "brokenwing-test-"); err != nil {
		logging.Fatal(err.Error())
	}

	testCode = path.Join(dir, "definitions.go")
	testFile = path.Join(dir, "definitions")

	err = ioutil.WriteFile(testCode, []byte(definitions.Code()), 0700|os.ModeTemporary)
	if err != nil {
		logging.Fatal(err.Error())
	}

	if forensic {
		logging.Noticef("Wrote '%s'. File will not be deleted, due to forensic mode.", string(testCode))
	} else {
		logging.Debugf("Wrote '%s'", string(testCode))
	}

	return
}
