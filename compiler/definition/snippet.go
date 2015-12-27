package definition

var snippet string = `
package main

import (
	. "github.com/dekelund/unbrokenwing/unbrokenwing"
	"github.com/dekelund/stdres"
	"strconv"
	"os"
	"log"
	"testing"
)

%s  // FIXME: This will not work with "import ("

func setup() {
%s
}

func main() {
	if pretty, err := strconv.ParseBool(os.Args[1]); err != nil {
		log.Fatal("Error configuring pretty print: ", err)
	} else if pretty {
		stdres.EnableColor()
	} else {
		stdres.DisableColor()
	}

	setup()
	feature := NewFeature(os.Stdin)
	suite := NewSuite()
	t := testing.T{}
	suite.Test(*feature, &t)
}`
