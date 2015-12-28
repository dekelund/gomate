// Package feature provides ...
package feature

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dekelund/unbrokenwing/global"
)

func getFeaturePaths(path string) (list []string) {
	dir, err := os.Open(path)

	if err != nil {
		global.Fatalf("Error opening input file:", err)
	}

	defer dir.Close()

	for names, err := dir.Readdirnames(10); err != io.EOF; names, err = dir.Readdirnames(10) {
		if err != nil {
			global.Fatalf("Error listing files:", err)
		}

		for _, name := range names {
			if !strings.HasSuffix(name, ".feature") {
				global.Debug(fmt.Sprintf("Ignoring non-feature file: '%s'\n", string(name)))
				continue
			}

			fpath, err := filepath.Abs(filepath.Join(path, name))

			if err != nil {
				global.Fatal(err.Error())
			}

			list = append(list, fpath)
		}
	}

	return list
}

func getDefinitonPaths(path string) (list []string) {
	if _, err := isDir(path); err != nil {
		global.Fatal(err.Error())
	}

	dir, err := os.Open(path)

	defer dir.Close()

	if err != nil {
		global.Fatalf("Error opening input file:", err)
	}

	for names, err := dir.Readdirnames(10); err != io.EOF; names, err = dir.Readdirnames(10) {
		if err != nil {
			global.Fatalf("Error listing files:", err)
		}

		for _, name := range names {
			if !strings.HasSuffix(name, ".go") {
				global.Debug(fmt.Sprintf("Ignoring non-definition file: '%s'\n", string(name)))
				continue
			}

			defPath, err := filepath.Abs(filepath.Join(path, name))

			if err != nil {
				global.Fatal(err.Error())
			}

			list = append(list, defPath)
		}
	}

	return list
}

// List represents the files and subdirectories files from a feature folder, including step definitions.
type List struct {
	Features    []string
	Definitions []string
}

// ParseDir make use of tools input data to generate definions binary and features struct.
// fpath represents a relative path, to a .feature file or a dir with .feature files.
// defPattern represents definitions folders name, shall be located in features directory.
// Function returns a list of features found in features file/dir and corresponding definitions.
// An error will be returned if error occur, if not caller are responsible to call Definitions.Remove().
func ParseDir(fpath string) (list List, err error) {
	var dir bool

	if fpath, err = filepath.Abs(fpath); err != nil {
		global.Fatal(err.Error())
	}

	global.Debug(fmt.Sprintf("Going to test %s\n", fpath))

	if dir, err = isDir(fpath); err != nil {
		return
	} else if dir {
		list.Features = getFeaturePaths(fpath)
	} else {
		list.Features = []string{fpath}
		fpath = filepath.Dir(fpath) // Point fpath to dir
	}

	list.Definitions = getDefinitonPaths(
		path.Join(fpath, global.Settings.DefPattern), //Dir with .go-definitions
	)

	return
}

func isDir(file string) (bool, error) {
	inputFile, err := os.Open(file)
	if err != nil {
		err := errors.New("Error opening input file: " + file)
		return false, err
	}

	defer inputFile.Close()

	info, err := inputFile.Stat()
	if err != nil {
		err := errors.New("Error for stat of input file: " + file)
		return false, err
	}

	return info.IsDir(), nil
}
