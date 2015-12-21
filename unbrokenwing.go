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
	"github.com/dekelund/unbrokenwing/global"
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

var CWD string = "."

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	CWD = cwd
}

func main() {
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
			Usage: "Relative path, to a feature-file or -directory (Current value: " + CWD + ").",
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

func setupGlobals(c *cli.Context) {
	global.Debug = c.GlobalBool("debug")
	global.PPrint = c.GlobalBool("pretty")
	global.DefPattern = c.GlobalString("step-definitions")

	global.CWD = CWD

	if global.PPrint {
		stdres.EnableColor()
	} else {
		stdres.DisableColor()
	}
}

func listFeatureFilesCmd(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	_, features := parseDir(dir, global.DefPattern, global.Debug)

	for i, feature := range features {
		path := CWD + PathSeparator
		fmt.Printf("\t%2d) %s\n", i, strings.TrimPrefix(feature, path))
	}
}

func listFeaturesCmd(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	_, features := parseDir(dir, global.DefPattern, global.Debug)

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

		path := CWD + PathSeparator
		fmt.Print("\n# ", strings.TrimPrefix(feature, path), "\n", text, "\n")
	}
}

func printDefinitionsCodeCmd(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	definitions, _ := parseDir(dir, global.DefPattern, global.Debug)

	fmt.Println(definitions.Code())
}

// testCmd search, compile and execute features defined in Gherik format where behaviours are defined in Go-Lang based files.
// Behaviours might be undefined, which will end up as red text in stdout if the context c has pretty print enabled.
func testCmd(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	definitions, features := parseDir(dir, global.DefPattern, global.Debug)

	if !global.Debug {
		defer definitions.Remove()
	}

	for _, file := range features {
		fd, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()

		definitions.Run(fd, global.PPrint, global.Debug)
	}
}

func parseDir(path, defPattern string, debug bool) (definition.Definitions, []string) {
	var err error
	var list = feature.List{}
	var defs = []io.Reader{}

	if list, err = feature.ParseDir(path); err != nil {
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
