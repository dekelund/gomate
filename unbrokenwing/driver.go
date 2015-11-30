package unbrokenwing

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"testing"

	"github.com/dekelund/stdres"
)

var buffer stdres.Buffer

func init() {
	buffer = stdres.Buffer{}
}

var featureRegexp = regexp.MustCompile("^Feature: (?P<name>.+)$")
var descriptionRegexp = regexp.MustCompile("^  ((?P<text>.+))?$")
var scenarioRegexp = regexp.MustCompile("^  Scenario: (?P<description>[a-zA-Z ]+)")
var stepRegexp = regexp.MustCompile("^    (?P<cmd>Given|When|Then|But|And) (?P<description>.+)$")
var emptyLineRexexp = regexp.MustCompile("^[\t ]+$")

// NewFeature scans FeatureFile for lines starting with
// "Feature:" followed by feature name, description
// and different scenarios. All scenarios including
// description are then returned as a Feature.
func NewFeature(reader io.Reader) (feature *Feature) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if featureRegexp.MatchString(line) {
			feature = scanFeature(getArgs(featureRegexp, line), scanner)
		} else {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
	}

	return
}

func scanFeature(regexpMap Args, scanner *bufio.Scanner) (feature *Feature) {
	feature = &Feature{}

	feature.Name = regexpMap["name"]
	feature.Description = ""

	for scanner.Scan() {
		line := scanner.Text()
		if scenarioRegexp.MatchString(line) {
			feature.Scenarios = append(feature.Scenarios, scanScenario(getArgs(scenarioRegexp, line), scanner))
		} else if descriptionRegexp.MatchString(line) {
			feature.Description += "\n" + line
		} else {
			feature.Description += "\n" + line
		}
	}

	return
}

func getArgs(re *regexp.Regexp, line string) (regexpMap Args) {
	regexpMap = Args{}

	if matches := re.FindStringSubmatch(line); matches != nil {
		for i, name := range re.SubexpNames() {
			regexpMap[name] = matches[i]
		}
	}

	return
}

func scanScenario(regexpMap Args, scanner *bufio.Scanner) (scenario Scenario) {
	scenario.Description = regexpMap["description"]

	for scanner.Scan() {
		line := scanner.Text()

		if emptyLineRexexp.MatchString(line) {
			//buffer.Println(fmt.Sprintf("WARNING: Can't parse step \"%s\"\n", line)).Result = stdres.FAILURE
			continue
		} else if stepRegexp.MatchString(line) {
			scenario.Steps = append(scenario.Steps, scanStep(getArgs(stepRegexp, line), scanner))
		} else {
			return
		}
	}

	return
}

func scanStep(regexpMap Args, scanner *bufio.Scanner) (step Step) {
	step.Description = regexpMap["description"]
	step.Cmd = regexpMap["cmd"]

	return
}

// Test runs a Feature and record test result.
// Test results are based on bahaviours
// supplied by one of following commands:
// Given, When, Then, But, And, Asterix.
func (ts *suite) Test(feature Feature, t *testing.T) error {
	err := ts.testFeature(feature, t)

	if err != nil {
		buffer.Println(fmt.Sprintf("Test framework failed for: %s\n", feature.Name)).Result = stdres.FAILURE
		t.Fail()
	}

	return err
}

func (ts *suite) testFeature(feature Feature, t *testing.T) error {
	ts.totalFeatures += 1

	featureText := buffer.Println(fmt.Sprintf("Feature: %s\n", feature.Name))
	featureText.Result = stdres.SUCCESS // Assume succes until something else has been proven
	defer func() {
		buffer.Println(fmt.Sprintf("%s", ts)).Result = stdres.PLAIN
		buffer.Println("\n    You can implement step definition for undefined steps with these snippets:").Result = stdres.PLAIN
		buffer.Println(fmt.Sprintf("%s", ts.Snippets())).Result = stdres.INFO
		buffer.Flush()
	}()

	for _, scenario := range feature.Scenarios {
		err := ts.testScenario(scenario)

		switch err.(type) {
		case nil:
			ts.successFeatures += 1
		case PendingError:
			ts.pendingFeatures += 1
			featureText.Result = stdres.PENDING
		case NotImplError:
			featureText.Result = stdres.UNKNOWN
		default:
			ts.failuresFeatures += 1
			featureText.Result = stdres.FAILURE
		}
	}

	return nil
}

func (ts *suite) testScenario(scenario Scenario) error {
	var notimplemented, pending, failure bool // TODO: We are not able to identify not implemented scenarios(?)
	var toReturn error

	scenarioText := buffer.Println(fmt.Sprintf("  Scenario: %s\n", scenario.Description))
	scenarioText.Result = stdres.UNKNOWN
	ts.totalScenarios += 1

	for _, step := range scenario.Steps {
		optout := notimplemented || pending || failure
		err := ts.testStep(step, optout)

		switch e := err.(type) {
		case nil:
			continue // Successfully executed, nothing to do
		case PendingError:
			pending = true
		case NotImplError:
			notimplemented = true
			ts.missingImpl[e.Snippet()] = true
		default:
			failure = true
		}

		switch e := toReturn.(type) {
		case nil:
			toReturn = err
		case NotImplError:
			toReturn = err
		case PendingError:
			switch e := toReturn.(type) {
			case NotImplError:
				_ = e
				toReturn = err
			}
		default:
			_ = e
			switch e := toReturn.(type) {
			case NotImplError, PendingError:
				_ = e
				toReturn = err
			}
			failure = true
		}
	}

	if failure {
		scenarioText.Result = stdres.FAILURE
		ts.failuresScenarios += 1
	} else if pending {
		scenarioText.Result = stdres.PENDING
		ts.pendingScenarios += 1
	} else if notimplemented { // TODO: We are not able to identify not implemented scenarios(?)
		scenarioText.Result = stdres.UNKNOWN
		ts.undefinedScenarios += 1
	} else {
		scenarioText.Result = stdres.SUCCESS
		ts.successScenarios += 1
	}

	return toReturn //NotImplemented(scenario)
}

func (ts *suite) testStep(step Step, optout bool) error {
	text := buffer.Println(fmt.Sprintf("    %s %s", step.Cmd, step.Description))
	text.Result = stdres.UNKNOWN
	defer func() {
		buffer.Println("").Result = stdres.INFO
	}()

	ts.totalSteps += 1

	if optout {
		defer func() {
			text.Result = stdres.FAILURE
		}()
	}

	for _, impl := range stepRegister {
		match, err := impl(step.Description, optout)

		if !match {
			continue
		}

		switch err.(type) {
		case nil:
			if optout {
				ts.optoutSteps += 1
			} else {
				ts.successSteps += 1
			}
		case PendingError:
			ts.pendingSteps += 1
			text.Result = stdres.PENDING
		default:
			ts.failuresSteps += 1
			text.Result = stdres.FAILURE
		}

		return err
	}

	ts.undefinedSteps += 1
	return NotImplemented(step)
}
