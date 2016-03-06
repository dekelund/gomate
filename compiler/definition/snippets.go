package definition

var test string = `
package main

import (
	. "gomate.io/gomate/unbrokenwing"
	"github.com/dekelund/stdres"
	"strconv"
	"os"
	"log"
	"testing"

	%s
)

func main() {
	if pretty, err := strconv.ParseBool(os.Args[1]); err != nil {
		log.Fatal("Error configuring pretty print: ", err)
	} else if pretty {
		stdres.EnableColor()
	} else {
		stdres.DisableColor()
	}

	feature := NewFeature(os.Stdin)
	suite := NewSuite()
	t := testing.T{}
	suite.Test(*feature, &t)
}`

var scaffold string = `
package main

import (
	. "gomate.io/gomate/unbrokenwing"
	"bufio"
	"fmt"
	"os"

	%s
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, err := reader.ReadString('\n')
		if err != nil {
			panic(err.Error())
		}

		errors := ExecuteCMD(cmd)

		for _, err := range errors {
			fmt.Println(err.Error())
		}

		if len(errors) > 0 {
			panic("To many errors, panic!")
		}
	}
}`
