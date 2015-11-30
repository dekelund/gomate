package definition_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dekelund2/unbrokenwing/compiler/definition"
)

func Example_NewDefinition() {

	buffer := bytes.NewBufferString(`
package step_definitions

Given("^I'm successfully logged in as an admin in users pane$", func(args Args) error {
	_ = ioutil.Discard
	return Pending("Not implemented")
})
And("^I fill in a new developer named hacker with password changeme$", func(args Args) error {
	_ = strings.Join([]string{}, "")
	return Pending("Not implemented")
})

import "strings"

	`)

	def := definition.NewDefinition(buffer)

	fmt.Println(def.Code())
	// Output:
	// package main
	//
	// import (
	// 	. "github.com/dekelund2/unbrokenwing/unbrokenwing"
	// 	"github.com/dekelund/stdres"
	// 	"strconv"
	// 	"os"
	// 	"log"
	// 	"testing"
	// )
	//
	// import "strings"  // FIXME: This will not work with "import ("
	//
	// func setup() {
	// Given("^I'm successfully logged in as an admin in users pane$", func(args Args) error {
	// 	_ = ioutil.Discard
	// 	return Pending("Not implemented")
	// })
	// And("^I fill in a new developer named hacker with password changeme$", func(args Args) error {
	// 	_ = strings.Join([]string{}, "")
	// 	return Pending("Not implemented")
	// })
	// }
	//
	// func main() {
	// 	file := os.Args[1]
	// 	fd, err := os.Open(file)  // Open *.feature file
	// 	if err != nil {
	// 		log.Fatal("Error opening input file:", err)
	// 	}
	//
	// 	defer fd.Close()
	//
	// 	if pretty, err := strconv.ParseBool(os.Args[2]); err != nil {
	// 		log.Fatal("Error configuring pretty print: ", err)
	// 	} else if pretty {
	// 		stdres.EnableColor()
	// 	} else {
	// 		stdres.DisableColor()
	// 	}
	//
	// 	setup()
	// 	feature := NewFeature((*FeatureFile)(fd))
	// 	suite := NewSuite()
	// 	t := testing.T{}
	// 	suite.Test(*feature, &t)
	// }

}

func Example_Run() {

	definitions := definition.NewDefinitions([]io.Reader{
		bytes.NewBufferString(`
package step_definitions

Given("^I'm successfully logged in as an admin in users pane$", func(args Args) error {
	_ = ioutil.Discard
	return Pending("Not implemented")
})
And("^I fill in a new developer named hacker with password changeme$", func(args Args) error {
	_ = strings.Join([]string{}, "")
	return Pending("Not implemented")
})

import "strings"

	`),
	}, false)

	features := bytes.NewBufferString(`
Feature: Manage users
    Administrators are able to manage all users.
    It's possible to add and remove users, but also
    change user kind, password and username.

    Following user kinds are supported by default
    setup:

    * Administrators:
        Setup new users and configures accessibility
        (Users and network interfaces)
    * Developers:
        Are able to setup new projects containing
        build and test configuration.
    * Applications:
        Rest and API systems, for instance external webpages.
    * Others:
        This category contains users with limited access.
        Project managers and field engineers are canonical users.
        Typical user scenario: tag release, download doployable
        build.

  Scenario: Create a new developer
    Given I'm successfully logged in as an admin in users pane
    And I fill in a new developer named hacker with password changeme
    When I press the create button
    Then only one user-record with name hacker should exist
    And user hacker should have password changeme
`)

	var pprint, debug bool

	definitions.Run(features, pprint, debug)

	// Output:
	// Feature: Manage users
	//
	//   Scenario: Create a new developer
	//
	//     Given I'm successfully logged in as an admin in users pane
	//
	//     And I fill in a new developer named hacker with password changeme
	//
	//     When I press the create button
	//
	//     Then only one user-record with name hacker should exist
	//
	//     And user hacker should have password changeme
	//
	//     1 scenario (0 undefined, 0 failures, 1 pending)
	//     5 steps (3 undefined, 0 failures, 1 pending, 1 optout)
	//
	//     You can implement step definition for undefined steps with these snippets:
	//
	//     And("^user hacker should have password changeme$", func(args Args) error {
	//         return Pending("Not implemented")
	//     })
	//
	//     Then("^only one user-record with name hacker should exist$", func(args Args) error {
	//         return Pending("Not implemented")
	//     })
	//
	//     When("^I press the create button$", func(args Args) error {
	//         return Pending("Not implemented")
	//     })
}
