package definition

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	. "gomate.io/gomate/global"
)

type Definition struct {
	pkgs    []string
	tmpDir  string // TODO: Remove tmpDir and command
	command string
}

// Run takes a DSL written feature from io.Reader, compiles behaviour code and supplies the feature code
// into executed behaviour code. After execution of the binary, the result are written to STDOUT. Method
// respects global.PPrint to enable/disable pretty print i.e., colors enabled.
func (definition Definition) Run(features io.Reader) {
	(&definition).compile()

	if !Settings.Forensic {
		defer (&definition).remove()
	}

	featureLines, err := ioutil.ReadAll(features)

	if err != nil {
		Fatal(err.Error())
	}

	gorun := exec.Command(definition.command, strconv.FormatBool(Settings.PPrint))
	if stdin, err := gorun.StdinPipe(); err != nil {
		Fatal(err.Error())
	} else if n, err := stdin.Write(featureLines); err != nil {
		Fatal(err.Error())
	} else if n != len(featureLines) {
		Fatal("Behaviour binary file was not able to read all defined features")
	} else if err := stdin.Close(); err != nil {
		Fatal(err.Error())
	}

	if output, err := gorun.CombinedOutput(); err != nil {
		Fatal(err.Error())
	} else {
		Info(string(output))
	}
}

func (definition *Definition) compile() {
	var err error
	var output []byte

	dir, testCode, testFile := definition.store()

	goimport := exec.Command("goimports", "-w=true", testCode)
	gofmt := exec.Command("go", "fmt", testCode)
	gobuild := exec.Command("go", "build", "-o", testFile, testCode)

	if err = goimport.Run(); err != nil {
		Err(err.Error())
	}

	if err = gofmt.Run(); err != nil {
		Err(err.Error())
	}

	if output, err = gobuild.CombinedOutput(); err != nil {
		Err(string(output))
	}

	definition.tmpDir = dir
	definition.command = testFile
}

func (definition Definition) store() (dir, testCode, testFile string) {
	var err error

	if dir, err = ioutil.TempDir("", "brokenwing-test-"); err != nil {
		Fatal(err.Error())
	}

	testCode = path.Join(dir, "definition.go")
	testFile = path.Join(dir, "definition")

	if ioutil.WriteFile(testCode, []byte(definition.TestCode()), 0700|os.ModeTemporary); err != nil {
		Fatal(err.Error())
	}

	if Settings.Forensic {
		Noticef("Wrote '%s'. File will not be deleted, due to forensic mode.", string(testCode))
	} else {
		Debugf("Wrote '%s'", string(testCode))
	}

	return
}

// Remove will delete temporary file containing the generated step definition
// After this method has been called, it's no longer possible to execute Run
// method.
func (definition *Definition) remove() {
	if err := os.RemoveAll(definition.tmpDir); err != nil {
		Err(err.Error())
		return
	}
}

func (definition Definition) imports() (imports string) {
	if len(definition.pkgs) == 0 {
		return
	}

	imports = "_ \"" + strings.Join(definition.pkgs, "\"\n_ \"") + "\"\n"
	return
}

// TestCode returns code used by gomate tool
// to run test suits based on packages provided
// when allocating Definition.
func (definition Definition) TestCode() string {
	return fmt.Sprintf(test, definition.imports())
}

// Scaffold returns code to be used as a baseline
// when implementing code making use of optional
// command patterns in definition files.
func (definition Definition) Scaffold() string {
	return fmt.Sprintf(scaffold, definition.imports())
}

// New allocates and returns definition struct
// that correspond to subpackages fond under import
// paths provided as argument.
//
// All packages provided in paths slice need to
// be grand child in the file structure to, and thereby
// relative to, GOPATH.
//
// Returned definition structure maps against
// definitions written in step_definition folders.
//
// Note: step-definition folder name might be
// changed based on Settings.DefPattern value.
func New(paths []string) (defs Definition) {
	pkgs := map[string]bool{}
	paths = PKGPathToPath(Settings.CWD, Settings.GOSRCPATH, paths)

	for _, path := range paths {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			dir, file := filepath.Split(path)

			if filepath.Ext(file) != ".go" {
				return nil
			}

			if filepath.Base(dir) == Settings.DefPattern {
				if pkg, err := filepath.Rel(Settings.GOSRCPATH, dir); err != nil {
					panic(err.Error())
				} else {
					pkgs[pkg] = true
				}
			}

			return nil
		})
	}

	for pkg := range pkgs {
		defs.pkgs = append(defs.pkgs, pkg)
	}

	(&defs).compile()
	return
}
