package main

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"log/syslog"

	"github.com/codegangsta/cli"
	"github.com/dekelund/stdres"

	"gomate.io/gomate/compiler/definition"
	"gomate.io/gomate/compiler/feature"
	. "gomate.io/gomate/global"
	"gomate.io/gomate/internal/highlighter"
)

const (
	PathSeparator = string(os.PathSeparator)
)

var CWD string = "."

func init() {
	var err error

	if CWD, err = os.Getwd(); err != nil {
		Fatal(err.Error())
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
			Name:    "scaffold",
			Aliases: []string{},
			Usage:   "Create code that initiate alternative protocol commands from step definitions.",
			Flags:   []cli.Flag{},
			Action:  scaffoldCMD,
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
	Settings.CWD = CWD

	Settings.SysLog.Active = c.GlobalBool("syslog")
	Settings.SysLog.UDP = c.GlobalBool("syslog-udp")
	Settings.SysLog.RAddr = c.GlobalString("syslog-raddr")
	Settings.SysLog.Tag = c.GlobalString("syslog-tag")
	Settings.SysLog.Priority = syslog.Priority(c.GlobalInt("priority"))

	Settings.PPrint = c.GlobalBool("pretty")
	Settings.Forensic = c.GlobalBool("forensic")
	Settings.DefPattern = c.GlobalString("step-definitions")

	if Settings.PPrint {
		stdres.EnableColor()
	} else {
		stdres.DisableColor()
	}

	ReconfigureLogger()
}

func listFeatureFilesCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	_, features := parseDir(dir)

	for i, feature := range features {
		path := CWD + PathSeparator
		Infof("\t%2d) %s\n", i, strings.TrimPrefix(feature, path))
	}
}

func listFeaturesCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	_, features := parseDir(dir)

	for _, feature := range features {
		fileReader, err := os.Open(feature)
		if err != nil {
			Fatal(err.Error())
		}

		bytes, err := ioutil.ReadAll(fileReader)
		if err != nil {
			Fatal(err.Error())
		}

		text := string(bytes)

		if Settings.PPrint {
			text = highlighter.Feature(text)
		}

		path := CWD + PathSeparator
		Infof("\n# %s\n%s\n", strings.TrimPrefix(feature, path), text)
	}
}

func printDefinitionsCodeCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	definitions, _ := parseDir(dir)

	defs := definitions.TestCode()

	if Settings.PPrint {
		defs = highlighter.Definition(defs)
	}

	Infof(defs)
}

// testCMD search, compile and execute features defined in Gherik format where behaviours are defined in Go-Lang based files.
// Behaviours might be undefined, which will end up as red text in stdout if the context c has pretty print enabled.
func testCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	definitions, features := parseDir(dir)

	if !Settings.Forensic {
		defer definitions.Remove()
	}

	for _, file := range features {
		fd, err := os.Open(file)
		if err != nil {
			Fatal(err.Error())
		}
		defer fd.Close()

		definitions.Run(fd)
	}
}

func scaffoldCMD(c *cli.Context) {
	setupGlobals(c)
	dir := c.GlobalString("dir")

	definitions, _ := parseDir(dir)

	defs := definitions.ScaffoldCode()

	if Settings.PPrint {
		defs = highlighter.Definition(defs)
	}

	Infof(defs)
}

func parseDir(path string) (definition.Definitions, []string) {
	var err error
	var list = feature.List{}
	var defs = []io.Reader{}

	if list, err = feature.ParseDir(path); err != nil {
		Fatal(err.Error())
	}

	for _, def := range list.Definitions {
		file, err := os.Open(def)
		if err != nil {
			Fatal(err.Error())
		}

		defs = append(defs, io.Reader(file))
		defer file.Close()
	}

	return definition.NewDefinitions(defs), list.Features
}
