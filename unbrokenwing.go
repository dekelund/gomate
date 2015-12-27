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
			Action:  listFeatureFilesCMD,
		},
		{
			Name:    "features",
			Aliases: []string{},
			Usage:   "List features to STDOUT",
			Flags:   []cli.Flag{},
			Action:  listFeaturesCMD,
		},
		{
			Name:    "definitions",
			Aliases: []string{"defs", "code"},
			Usage:   "List behaviours to STDOUT",
			Flags:   []cli.Flag{},
			Action:  printDefinitionsCodeCMD,
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Tests either a test directory with features in it, or a .feature file",
			Flags:   []cli.Flag{},
			Action:  testCMD,
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

func listFeatureFilesCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	_, features := parseDir(dir)

	for i, feature := range features {
		path := CWD + PathSeparator
		fmt.Printf("\t%2d) %s\n", i, strings.TrimPrefix(feature, path))
	}
}

func listFeaturesCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	_, features := parseDir(dir)

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

		if global.PPrint {
			text = strings.Replace(text, "Feature: ", red+"Feature: "+reset, -1)
			text = strings.Replace(text, "Scenario: ", red+"Scenario: "+reset, -1)
			text = strings.Replace(text, " Given ", green+" Given "+reset, -1)
			text = strings.Replace(text, " And ", green+" And "+reset, -1)
			text = strings.Replace(text, " When ", blue+" When "+reset, -1)
			text = strings.Replace(text, " Then ", yellow+" Then "+reset, -1)
		}

		path := CWD + PathSeparator
		fmt.Print("\n# ", strings.TrimPrefix(feature, path), "\n", text, "\n")
	}
}

func printDefinitionsCodeCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	definitions, _ := parseDir(dir)

	text := string(definitions.Code())

	if global.PPrint {
		text = strings.Replace(text, "package main", yellow+"package "+blue+"main"+reset, -1)
		text = strings.Replace(text, "import ", yellow+"import"+reset+" ", -1)
		text = strings.Replace(text, "func ", blue+"func"+reset+" ", -1)
		text = strings.Replace(text, "func(", blue+"func"+reset+"(", -1)
		text = strings.Replace(text, "defer ", blue+"defer"+reset+" ", -1)
		text = strings.Replace(text, " error ", " "+red+"error"+reset+" ", -1)

		text = strings.Replace(text, "\nFeature(", red+"\nFeature"+reset+"(", -1)
		text = strings.Replace(text, "\nScenario(", red+"\nScenario"+reset+"(", -1)
		text = strings.Replace(text, "\nGiven(", green+"\nGiven"+reset+"(", -1)
		text = strings.Replace(text, "\nAnd(", green+"\nAnd"+reset+"(", -1)
		text = strings.Replace(text, "\nWhen(", blue+"\nWhen"+reset+"(", -1)
		text = strings.Replace(text, "\nThen(", yellow+"\nThen"+reset+"(", -1)

		text = strings.Replace(text, "main()", yellow+"main"+reset+"()", -1)
		text = strings.Replace(text, "setup()", yellow+"setup"+reset+"()", -1)

		text = strings.Replace(text, " Pending(", " "+yellow+"Pending"+reset+"(", -1)
		text = strings.Replace(text, "os.Open(", yellow+"os.Open"+reset+"(", -1)
		text = strings.Replace(text, "fd.Close(", yellow+"fd.Close"+reset+"(", -1)
		text = strings.Replace(text, "suite.Test(", yellow+"suite.Test"+reset+"(", -1)
		text = strings.Replace(text, "ParseBool(", yellow+"ParseBool"+reset+"(", -1)
		text = strings.Replace(text, "stdres.EnableColor(", yellow+"stdres.EnableColor"+reset+"(", -1)
		text = strings.Replace(text, "stdres.DisableColor(", yellow+"stdres.DisableColor"+reset+"(", -1)
		text = strings.Replace(text, "os.Args[", yellow+"os.Args"+reset+"[", -1)
	}

	fmt.Println(text)
}

// testCMD search, compile and execute features defined in Gherik format where behaviours are defined in Go-Lang based files.
// Behaviours might be undefined, which will end up as red text in stdout if the context c has pretty print enabled.
func testCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	definitions, features := parseDir(dir)

	if !global.Debug {
		defer definitions.Remove()
	}

	for _, file := range features {
		fd, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()

		definitions.Run(fd)
	}
}

func parseDir(path string) (definition.Definitions, []string) {
	var err error
	var list = feature.List{}
	var defs = []io.Reader{}

	if list, err = feature.ParseDir(path); err != nil {
		log.Fatal(err)
	}

	for _, def := range list.Definitions {
		file, err := os.Open(def)
		if err != nil {
			log.Fatal(err)
		}

		defs = append(defs, io.Reader(file))
		defer file.Close()
	}

	return definition.NewDefinitions(defs), list.Features
}
