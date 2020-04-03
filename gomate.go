package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"log/syslog"

	"github.com/dekelund/stdres"
	"github.com/urfave/cli"

	"gomate.io/gomate/compiler/definition"
	"gomate.io/gomate/compiler/feature"
	"gomate.io/gomate/internal/highlighter"
	"gomate.io/gomate/logging"
)

const (
	pathSeparator = string(os.PathSeparator)
)

var settings struct {
	SysLog     logging.Settings
	Forensic   bool
	PPrint     bool
	CWD        string
	DefPattern string
}

var cwd = "."

func init() {
	settings.SysLog.Priority = syslog.LOG_INFO

	var err error

	if cwd, err = os.Getwd(); err != nil {
		logging.Fatal(err.Error())
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "gomate"
	app.Version = "0.1"
	app.Usage = "Run behaviour driven tests as Gherik features"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "syslog",
			Usage: "Redirect STDOUT to SysLog server",
		},
		cli.BoolFlag{
			Name:  "syslog-udp",
			Usage: "Use UDP instead of TCP",
		},
		cli.StringFlag{
			Name:  "syslog-raddr",
			Usage: "HOST/IP address to SysLog server",
			Value: "localhost",
		},
		cli.StringFlag{
			Name:  "syslog-tag",
			Usage: "Tag output with specified text string",
			Value: "gomate",
		},
		cli.IntFlag{
			Name: "priority",
			Usage: "Log priority, use bitwised values from /usr/include/sys/syslog.h e.g.," +
				" LOG_EMERG=" + strconv.Itoa(int(syslog.LOG_EMERG)) +
				" LOG_ALERT=" + strconv.Itoa(int(syslog.LOG_ALERT)) +
				" LOG_CRIT=" + strconv.Itoa(int(syslog.LOG_CRIT)) +
				" LOG_ERR=" + strconv.Itoa(int(syslog.LOG_ERR)) +
				" LOG_WARNING=" + strconv.Itoa(int(syslog.LOG_WARNING)) +
				" LOG_NOTICE=" + strconv.Itoa(int(syslog.LOG_NOTICE)) +
				" LOG_INFO=" + strconv.Itoa(int(syslog.LOG_INFO)) +
				" LOG_DEBUG=" + strconv.Itoa(int(syslog.LOG_DEBUG)),
			Value: int(syslog.LOG_INFO),
		},
		cli.BoolFlag{
			Name:  "pretty",
			Usage: "Print colorised result to STDOUT/STDERR",
		},
		cli.BoolFlag{
			Name:  "forensic",
			Usage: "A kind of development mode, all generated files will be kept",
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

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("exiting due to unexpected error: %s", err)
	}
}

func setupGlobals(c *cli.Context) {
	settings.CWD = cwd

	settings.SysLog.Active = c.GlobalBool("syslog")
	settings.SysLog.UDP = c.GlobalBool("syslog-udp")
	settings.SysLog.RAddr = c.GlobalString("syslog-raddr")
	settings.SysLog.Tag = c.GlobalString("syslog-tag")
	settings.SysLog.Priority = syslog.Priority(c.GlobalInt("priority"))

	settings.PPrint = c.GlobalBool("pretty")
	settings.Forensic = c.GlobalBool("forensic")
	settings.DefPattern = c.GlobalString("step-definitions")

	if settings.PPrint {
		stdres.EnableColor()
	} else {
		stdres.DisableColor()
	}

	logging.ReconfigureLogger(settings.SysLog)
}

func listFeatureFilesCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	_, features := parseDir(dir)

	for i, feature := range features {
		path := cwd + pathSeparator
		logging.Infof("\t%2d) %s\n", i, strings.TrimPrefix(feature, path))
	}
}

func listFeaturesCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	_, features := parseDir(dir)

	// #nosec
	for _, feature := range features {
		fileReader, err := os.Open(feature)
		if err != nil {
			logging.Fatal(err.Error())
		}

		bytes, err := ioutil.ReadAll(fileReader)
		if err != nil {
			logging.Fatal(err.Error())
		}

		text := string(bytes)

		if settings.PPrint {
			text = highlighter.Feature(text)
		}

		path := cwd + pathSeparator
		logging.Infof("\n# %s\n%s\n", strings.TrimPrefix(feature, path), text)
	}
}

func printDefinitionsCodeCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	definitions, _ := parseDir(dir)

	defs := definitions.Code()

	if settings.PPrint {
		defs = highlighter.Definition(defs)
	}

	logging.Infof(defs)
}

// testCMD search, compile and execute features defined in Gherik format where behaviours are defined in Go-Lang based files.
// Behaviours might be undefined, which will end up as red text in stdout if the context c has pretty print enabled.
func testCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	definitions, features := parseDir(dir)

	if !settings.Forensic {
		defer definitions.Remove()
	}

	// #nosec
	for _, file := range features {
		fd, err := os.Open(file)
		if err != nil {
			logging.Fatal(err.Error())
		}
		defer fd.Close()

		definitions.Run(fd, settings.PPrint)
	}
}

func parseDir(path string) (definition.Definitions, []string) {
	var err error
	var list = feature.List{}
	var defs = []io.Reader{}

	if list, err = feature.ParseDir(path, settings.DefPattern); err != nil {
		logging.Fatal(err.Error())
	}

	// #nosec
	for _, def := range list.Definitions {
		file, err := os.Open(def)
		if err != nil {
			logging.Fatal(err.Error())
		}

		defs = append(defs, io.Reader(file))
		defer file.Close()
	}

	return definition.NewDefinitions(defs, settings.Forensic), list.Features
}
