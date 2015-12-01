// Package feature provides ...
package feature

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getFeaturePaths(path string, debug bool) (list []string) {
	dir, err := os.Open(path)

	defer dir.Close()

	if err != nil {
		log.Fatal("Error opening input file:", err)
	}

	for names, err := dir.Readdirnames(10); err != io.EOF; names, err = dir.Readdirnames(10) {
		if err != nil {
			log.Fatal("Error listing files:", err)
		}

		for _, name := range names {
			if !strings.HasSuffix(name, ".feature") {
				if debug {
					fmt.Printf("Ignoring non-feature file: '%s'\n", string(name))
				}
				continue
			}

			fpath, err := filepath.Abs(filepath.Join(path, name))

			if err != nil {
				log.Fatal(err)
			}

			list = append(list, fpath)
		}
	}

	return list
}

func getDefinitonPaths(path string, debug bool) (list []string) {
	if _, err := isDir(path); err != nil {
		log.Fatal(err)
	}

	dir, err := os.Open(path)

	defer dir.Close()

	if err != nil {
		log.Fatal("Error opening input file:", err)
	}

	for names, err := dir.Readdirnames(10); err != io.EOF; names, err = dir.Readdirnames(10) {
		if err != nil {
			log.Fatal("Error listing files:", err)
		}

		for _, name := range names {
			if !strings.HasSuffix(name, ".go") {
				if debug {
					fmt.Printf("Ignoring non-feature file: '%s'\n", string(name))
				}
				continue
			}

			defPath, err := filepath.Abs(filepath.Join(path, name))

			if err != nil {
				log.Fatal(err)
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
	DefDir      string
}

// ParseDir make use of tools input data to generate definions binary and features struct.
// fpath represents a relative path, to a .feature file or a dir with .feature files.
// defPattern represents definitions folders name, shall be located in features directory.
// Function returns a list of features found in features file/dir and corresponding definitions.
// An error will be returned if error occur, if not caller are responsible to call Definitions.Remove().
func ParseDir(fpath, defPattern string, debug bool) (list List, err error) {
	var yes bool

	if fpath, err = filepath.Abs(fpath); err != nil {
		log.Fatal(err)
	}

	if debug {
		fmt.Printf("Going to test %s\n", fpath)
	}

	if yes, err = isDir(fpath); err != nil {
		return
	} else if yes {
		list.DefDir = path.Join(fpath, defPattern)
		list.Features = getFeaturePaths(fpath, debug)
	} else {
		list.DefDir = path.Join(filepath.Dir(fpath), defPattern)
		list.Features = []string{fpath}
	}

	list.Definitions = getDefinitonPaths(fpath, debug)
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
