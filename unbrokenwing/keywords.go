package unbrokenwing

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var didntMatch error = errors.New("Didn't match")

type cmd func(string) (int, interface{}, error)

var stepRegister = []func(string, bool) (match bool, err error){}
var cmdRegister = []cmd{}

// Execute none or many matching commands
// based on incomming command and step
// definitions.
func ExecuteCMD(cmd string) (cids []int, returns []interface{}, errors []error) {
	for _, definition := range cmdRegister {
		id, r, err := definition(cmd)

		if err == didntMatch {
			continue
		}

		cids = append(cids, id)
		errors = append(errors, err)
		returns = append(returns, r)
	}

	return
}

// https://github.com/cucumber/cucumber/wiki/Given-When-Then
func stepImplementation(step string, do func(Args) (interface{}, error), commands []string) {
	for _, cmdDef := range commands {
		r, err := regexp.Compile(cmdDef)

		if err != nil {
			fmt.Printf("WARNING: %s", err)
			continue
		}

		cmdRegister = append(cmdRegister, func(cmd string) (int, interface{}, error) {
			var id int = -1

			if r.MatchString(cmd) {
				args := getArgs(r, cmd)

				if cid, ok := args["gomateCMDId"]; ok {
					if i, err := strconv.Atoi(cid); err == nil {
						id = i
					}
				}

				result, err := do(args)
				return id, result, err
			}
			return id, nil, didntMatch
		})
	}

	stepRegister = append(stepRegister, func(line string, optout bool) (match bool, err error) {
		match = false
		r, err := regexp.Compile(step) // TODO reuse r from command-loop

		if err == nil {
			if r.MatchString(line) {
				match = true
				if !optout {
					_, err = do(getArgs(r, line))
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
func Given(step string, do func(Args) (interface{}, error), commands ...string) (err error) {
	stepImplementation(step, do, commands)
	return
}

// When are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
// As a last alternative regular expressions matching
// commands for implementation may be applied.
func When(step string, do func(Args) (interface{}, error), commands ...string) (err error) {
	stepImplementation(step, do, commands)
	return
}

// Then are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
// As a last alternative regular expressions matching
// commands for implementation may be applied.
func Then(step string, do func(Args) (interface{}, error), commands ...string) (err error) {
	stepImplementation(step, do, commands)
	return
}

// But are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
// As a last alternative regular expressions matching
// commands for implementation may be applied.
func But(step string, do func(Args) (interface{}, error), commands ...string) (err error) {
	stepImplementation(step, do, commands)
	return
}

// And are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
// As a last alternative regular expressions matching
// commands for implementation may be applied.
func And(step string, do func(Args) (interface{}, error), commands ...string) (err error) {
	stepImplementation(step, do, commands)
	return
}
