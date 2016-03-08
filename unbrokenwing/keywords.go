package unbrokenwing

import (
	"fmt"
	"regexp"
)

type cmd func(string) error

var stepRegister = []func(string, bool) (match bool, err error){}
var cmdRegister = []cmd{}

// Execute none or many matching commands
// based on incomming command and step
// definitions.
func ExecuteCMD(cmd string) (errors []error) {
	for _, definition := range cmdRegister {
		if err := definition(cmd); err != nil {
			errors = append(errors, err)
		}
	}

	return
}

// https://github.com/cucumber/cucumber/wiki/Given-When-Then
func stepImplementation(step string, do func(Args) error, commands []string) {
	for _, cmdDef := range commands {
		r, err := regexp.Compile(cmdDef)

		if err != nil {
			fmt.Printf("WARNING: %s", err)
			continue
		}

		cmdRegister = append(cmdRegister, func(cmd string) error {
			if r.MatchString(cmd) {
				return do(getArgs(r, cmd))
			}
			return nil
		})
	}

	stepRegister = append(stepRegister, func(line string, optout bool) (match bool, err error) {
		match = false
		r, err := regexp.Compile(step) // TODO reuse r from command-loop

		if err == nil {
			if r.MatchString(line) {
				match = true
				if !optout {
					err = do(getArgs(r, line))
				}
			}
		}

		return
	})
}

// Given are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
// As a last alternative regular expressions matching
// commands for implementation may be applied.
func Given(step string, do func(Args) error, commands ...string) (err error) {
	stepImplementation(step, do, commands)
	return
}

// When are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
// As a last alternative regular expressions matching
// commands for implementation may be applied.
func When(step string, do func(Args) error, commands ...string) (err error) {
	stepImplementation(step, do, commands)
	return
}

// Then are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
// As a last alternative regular expressions matching
// commands for implementation may be applied.
func Then(step string, do func(Args) error, commands ...string) (err error) {
	stepImplementation(step, do, commands)
	return
}

// But are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
// As a last alternative regular expressions matching
// commands for implementation may be applied.
func But(step string, do func(Args) error, commands ...string) (err error) {
	stepImplementation(step, do, commands)
	return
}

// And are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
// As a last alternative regular expressions matching
// commands for implementation may be applied.
func And(step string, do func(Args) error, commands ...string) (err error) {
	stepImplementation(step, do, commands)
	return
}
