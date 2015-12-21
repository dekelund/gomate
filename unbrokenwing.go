package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/dekelund/stdres"
	"github.com/dekelund/unbrokenwing/compiler/definition"
	"github.com/dekelund/unbrokenwing/compiler/feature"
)

// Foreground colors
const (
	reset   = "\033[00m"
	black   = "\033[30m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
)

const (
	PathSeparator = string(os.PathSeparator)
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	app := cli.NewApp()
	app.Name = "unbrokenwing"
	app.Usage = "Run behaviour driven tests as Gherik features"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "pretty",
			Usage: "Print colorised result to STDOUT/STDERR",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Verbose printing (Generated files will be kept)",
		},
		cli.StringFlag{
			Name:  "step-definitions",
			Value: "step_definitions",
			Usage: "Definitions folder name, should be located in features folder",
		},
		cli.StringFlag{
			Name:  "dir",
			Value: ".",
			Usage: "Relative path, to a feature-file or -directory (Current value: " + cwd + ").",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "feature-files",
			Aliases: []string{},
			Usage:   "List feature files to STDOUT",
			Flags:   []cli.Flag{},
			Action:  listFeatureFilesCmd,
		},
		{
			Name:    "features",
			Aliases: []string{},
			Usage:   "List features to STDOUT",
			Flags:   []cli.Flag{},
			Action:  listFeaturesCmd,
		},
		{
			Name:    "definitions",
			Aliases: []string{"defs", "code"},
			Usage:   "List behaviours to STDOUT",
			Flags:   []cli.Flag{},
			Action:  printDefinitionsCodeCmd,
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Tests either a test directory with features in it, or a .feature file",
			Flags:   []cli.Flag{},
			Action:  testCmd,
		},
	}

	app.Run(os.Args)
}

func listFeatureFilesCmd(c *cli.Context) {
	var debug bool = c.GlobalBool("debug")
	var defPattern string = c.GlobalString("step-definitions")
	var path string = c.GlobalString("dir")
	var features []string

	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	_, features = parseDir(path, defPattern, debug)

	for i, feature := range features {
		fmt.Printf("\t%2d) %s\n", i, strings.TrimPrefix(feature, cwd+PathSeparator))
	}
}

func listFeaturesCmd(c *cli.Context) {
	var debug bool = c.GlobalBool("debug")
	var defPattern string = c.GlobalString("step-definitions")
	var path string = c.GlobalString("dir")
	var features []string

	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	_, features = parseDir(path, defPattern, debug)

	for _, feature := range features {
		fileReader, err := os.Open(feature)
		if err != nil {
			log.Fatal(err)
		}

		bytes, err := ioutil.ReadAll(fileReader)
		if err != nil {
			log.Fatal(err)
		}

		text := string(bytes)
		text = strings.Replace(text, "Feature: ", red+"Feature: "+reset, -1)
		text = strings.Replace(text, "Scenario: ", red+"Scenario: "+reset, -1)
		text = strings.Replace(text, " Given ", green+" Given "+reset, -1)
		text = strings.Replace(text, " And ", green+" And "+reset, -1)
		text = strings.Replace(text, " When ", blue+" When "+reset, -1)
		text = strings.Replace(text, " Then ", yellow+" Then "+reset, -1)

		fmt.Print("\n# ", strings.TrimPrefix(feature, cwd+PathSeparator), "\n", text, "\n")
	}
}

func printDefinitionsCodeCmd(c *cli.Context) {
	var debug bool = c.GlobalBool("debug")
	var defPattern string = c.GlobalString("step-definitions")
	var path string = c.GlobalString("dir")
	var definitions definition.Definitions

	definitions, _ = parseDir(path, defPattern, debug)

	fmt.Println(definitions.Code())
}

// testCmd search, compile and execute features defined in Gherik format where behaviours are defined in Go-Lang based files.
// Behaviours might be undefined, which will end up as red text in stdout if the context c has pretty print enabled.
func testCmd(c *cli.Context) {
	var debug bool = c.GlobalBool("debug")
	var pretty bool = c.GlobalBool("pretty")
	var defPattern string = c.GlobalString("step-definitions")

	var path string = c.GlobalString("dir")

	if pretty {
		stdres.EnableColor()
	} else {
		stdres.DisableColor()
	}

	definitions, features := parseDir(path, defPattern, debug)

	if !debug {
		defer definitions.Remove()
	}

	for _, file := range features {
		fd, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()

		definitions.Run(fd, pretty, debug)
	}
}

func parseDir(path, defPattern string, debug bool) (definition.Definitions, []string) {
	var err error
	var list = feature.List{}
	var defs = []io.Reader{}

	if list, err = feature.ParseDir(path, defPattern, debug); err != nil {
		log.Fatal(err)
	}

	defFiles := []io.ReadCloser{}

	defer func() {
		for _, file := range defFiles {
			file.Close()
		}
	}() // Make sure to close all open files

	for _, def := range list.Definitions {
		if file, err := os.Open(def); err == nil {
			defFiles = append(defFiles, file)
		} else {
			log.Fatal(err)
		}
	}

	// FIXME figure out if it's possible to ignore convert
	for _, def := range defFiles {
		defs = append(defs, io.Reader(def))
	}

	return definition.NewDefinitions(defs, debug), list.Features
}
