package main

import (
	"io"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/dekelund/stdres"
	"github.com/dekelund/unbrokenwing/compiler/definition"
	"github.com/dekelund/unbrokenwing/compiler/feature"
)

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
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	app.Commands = []cli.Command{
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Tests either a test directory with features in it, or a .feature file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "feature",
					Value: ".",
					Usage: "Specifies file or directory to test, defaults to current working directory. (Right now: " + cwd + ")",
				},
			},
			Action: testCmd,
		},
	}

	app.Run(os.Args)
}

// testCmd search, compile and execute features defined in Gherik format where behaviours are defined in Go-Lang based files.
// Behaviours might be undefined, which will end up as red text in stdout if the context c has pretty print enabled.
func testCmd(c *cli.Context) {
	var debug bool = c.GlobalBool("debug")
	var pretty bool = c.GlobalBool("pretty")
	var defPattern string = c.GlobalString("step-definitions")

	var path string = c.String("feature")

	if pretty {
		stdres.EnableColor()
	} else {
		stdres.DisableColor()
	}

	fdList, err := feature.ParseDir(path, defPattern, debug)
	if err != nil {
		log.Fatal(err)
	}

	openDefList := []io.ReadCloser{}

	defer func() {
		for _, file := range openDefList {
			file.Close()
		}
	}() // Make sure to close all open files

	for _, def := range fdList.Definitions {
		if file, err := os.Open(def); err == nil {
			openDefList = append(openDefList, file)
		} else {
			log.Fatal(err)
		}
	}

	// Ugly convert for now
	openDefListConv := []io.Reader{}
	for _, def := range openDefList {
		openDefListConv = append(openDefListConv, io.Reader(def))
	}

	definitions := definition.NewDefinitions(openDefListConv, debug)

	//if err != nil {
	//	log.Panic(err)
	//}

	//dir, err = compiler.ParseDir(path, defPattern, debug)

	//if err != nil {
	//log.Panic(err)
	//}

	if !debug {
		defer definitions.Remove()
	}

	for _, file := range fdList.Features {
		fd, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()

		definitions.Run(fd, pretty, debug)
	}
}
