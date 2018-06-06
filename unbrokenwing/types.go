package unbrokenwing

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

// Args provides data structure for arguments
// supplied to the step definition.
type Args map[string]string

// Step corresponds to the a function related to
// a Given, When and Then-step.
//
// Cmd correspons to one of following commands: Given, When, Then, But, And
// Description contains the rest of the text that follows after the command.
type Step struct {
	Cmd         string
	Description string
}

// String returns the original text before broken down to cmd and description.
func (step Step) String() string {
	return fmt.Sprintf("%s %s", step.Cmd, step.Description)
}

// Scenario contains data structure matching scenarios in Gherkin.
// Description holds all text from scenario line till first scenario step.
type Scenario struct {
	Description string
	Steps       []Step
}

func (scenario Scenario) String() string {
	return fmt.Sprintf("Scenario: %s\n", scenario.Description)
}

// Feature contains data structure matching features in Gherkin.
// Each Scenario in Scenarios contains Description and scenario
// steps according to Gherkin scenarios.
type Feature struct {
	Name        string
	Description string
	Scenarios   []Scenario
}

func (feature Feature) String() string {
	return fmt.Sprintf("Feature: %s\n%s\n", feature.Name, feature.Description)
}

// Suite interface provides measures to run feature and record test result.
// Test results are based on bahaviours supplied by one of following commands:
// Given, When, Then, But, And, Asterix.
// String function returns test result as string, suitable to be printed to stdout.
type Suite interface {
	String() string
	Test(feature Feature, t *testing.T) error
}

// NewSuite generates built-in Suite implementation.
func NewSuite() Suite {
	s := suite{}
	s.missingImpl = map[string]bool{}

	return &s
}

type suite struct {
	totalFeatures  int
	totalScenarios int
	totalSteps     int

	undefinedScenarios int
	undefinedSteps     int

	successFeatures  int
	successScenarios int
	successSteps     int
	optoutSteps      int // Step not executed due to failure in earlier execution

	failuresFeatures  int
	failuresScenarios int
	failuresSteps     int

	pendingFeatures  int
	pendingScenarios int
	pendingSteps     int

	missingImpl map[string]bool
}

type byKey []string

func (a byKey) Len() int           { return len(a) }
func (a byKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byKey) Less(i, j int) bool { return a[i] < a[j] }

// snippets generates behaviour snippets based on Gherkin scenario steps.
func (ts suite) snippets() string {
	keys := make([]string, 0, len(ts.missingImpl))

	for k := range ts.missingImpl {
		keys = append(keys, k)
	}

	sort.Sort(byKey(keys))

	return strings.Join(keys, "\n")
}

// String function returns test result as string, suitable to be printed to stdout.
func (ts suite) String() string {
	return fmt.Sprintf("    %d scenario (%d undefined, %d failures, %d pending)\n    %d steps (%d undefined, %d failures, %d pending, %d optout)",
		ts.totalScenarios, ts.undefinedScenarios, ts.failuresScenarios, ts.pendingScenarios,
		ts.totalSteps, ts.undefinedSteps, ts.failuresSteps, ts.pendingSteps, ts.optoutSteps,
	)
}
