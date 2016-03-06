package unbrokenwing

import "regexp"

var stepRegister = []func(string, bool) (match bool, err error){}

// https://github.com/cucumber/cucumber/wiki/Given-When-Then
func stepImplementation(step string, do func(Args) error) {
	stepRegister = append(stepRegister, func(line string, optout bool) (match bool, err error) {
		match = false
		r, err := regexp.Compile(step)

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
func Given(step string, do func(Args) error) (err error) { stepImplementation(step, do); return }

// When are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
func When(step string, do func(Args) error) (err error) { stepImplementation(step, do); return }

// Then are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
func Then(step string, do func(Args) error) (err error) { stepImplementation(step, do); return }

// But are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
func But(step string, do func(Args) error) (err error) { stepImplementation(step, do); return }

// And are used to map scenario steps with behaviours,
// this is mapped by matching regular expression in first
// argument against scenario step in Gherkin language.
func And(step string, do func(Args) error) (err error) { stepImplementation(step, do); return }
