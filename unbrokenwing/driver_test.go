package unbrokenwing_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/dekelund/stdres"
	. "gomate.io/gomate/unbrokenwing"
)

func ExampleSuite_Test() {
	stdres.DisableColor()

	Given("^I'm successfully logged in as an admin in users pane$", func(args Args) error {
		_ = ioutil.Discard
		return Pending("Not implemented")
	})
	And("^I fill in a new developer named hacker with password changeme$", func(args Args) error {
		_ = strings.Join([]string{}, "")
		return Pending("Not implemented")
	})

	buffer := bytes.NewBufferString(`
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

	feature := NewFeature(buffer)
	suite := NewSuite()
	t := testing.T{}
	suite.Test(*feature, &t)

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
