package definition

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var emptyLineRexexp = regexp.MustCompile("^[\t ]*$")

// Definition represents a parsed step definition, typically located in step_definition folder underneath features folder.
type Definition struct {
	imports []string
	funcs   []string
}

type stepDefinitions []Definition

type Definitions struct {
	defs    stepDefinitions
	tmpDir  string
	command string
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
// After execution of the binary, the result are written to STDOUT. Caller are able to specify pretty
// print i.e., colors enabled. It's also possible to enable/disable debug output via method argument.
// NOTE: It's callers responsibility to call Remove method before Definitions instance are garbage
// collected.
func (definitions Definitions) Run(features io.Reader, pprint, debug bool) {
	var dir string
	var featureLines []byte
	var err error
	var output []byte

	if featureLines, err = ioutil.ReadAll(features); err != nil {
		log.Panic(err.Error())
	}

	if dir, err = ioutil.TempDir("", "brokenwing-test-"); err != nil {
		log.Panic(err.Error())
	}

	featureFile := path.Join(dir, "features")
	if ioutil.WriteFile(featureFile, featureLines, 0700|os.ModeTemporary); err != nil {
		log.Panic(err.Error())
	}

	gorun := exec.Command(definitions.command, featureFile, strconv.FormatBool(pprint))

	if output, err = gorun.CombinedOutput(); err != nil {
		fmt.Println(string(output))
	} else {
		fmt.Println(string(output))
	}
}

// NewDefinition reads and parse "in" assuming that it contains content from a step definition file.
// Lines defining package names are ommited from resulting Definition instance.
func NewDefinition(in io.Reader) Definition {
	code, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatal(err)
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
// Lines defining package names are ommited from resulting Definition instance.
// Caller has the possibility to enable debug output via method argument.
func NewDefinitions(defs []io.Reader, debug bool) Definitions {
	definitions := stepDefinitions{}

	for _, reader := range defs {
		definitions = append(definitions, NewDefinition(reader))
	}

	dir, cmd := definitions.compile(debug)
	return Definitions{definitions, dir, cmd}
}

// Remove will remove temporary file containing the generated step definition
// After this method has been called, it's no longer possible to execute Run
// method.
func (definitions Definitions) Remove() {
	os.Remove(definitions.tmpDir)
}

func (definitions stepDefinitions) compile(debug bool) (string, string) {
	var err error
	var output []byte

	dir, testCode, testFile := definitions.store(debug)

	goimport := exec.Command("goimports", "-w=true", testCode)
	gofmt := exec.Command("go", "fmt", testCode)
	gobuild := exec.Command("go", "build", "-o", testFile, testCode)

	if err = goimport.Run(); err != nil {
		fmt.Println(err.Error())
	}

	if err = gofmt.Run(); err != nil {
		fmt.Println(err.Error())
	}

	if output, err = gobuild.CombinedOutput(); err != nil {
		fmt.Println(string(output))
	}

	return dir, testFile
}

func (definitions stepDefinitions) store(debug bool) (dir, testCode, testFile string) {
	var err error

	if dir, err = ioutil.TempDir("", "brokenwing-test-"); err != nil {
		log.Panic(err.Error())
	}

	testCode = path.Join(dir, "definitions.go")
	testFile = path.Join(dir, "definitions")

	if ioutil.WriteFile(testCode, []byte(definitions.Code()), 0700|os.ModeTemporary); err != nil {
		log.Panic(err.Error())
	}

	if debug {
		fmt.Println("Wrote '", string(testCode), "'. File will not be deleted, due to debug mode.")
	}

	return
}
