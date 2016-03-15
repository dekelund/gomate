package unbrokenwing

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
)

const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32602
)

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type RPCMessage struct {
	// ID might be int, string and NULL according to JSON-RPC 2.0 , but we assume integer value
	// TODO: Remove string assumption from values in Params
	ID      int               `json:"id"`
	JSONRPC string            `json:"jsonrpc"` // Should be "2.0"
	Method  string            `json:"method"`
	Params  map[string]string `json:"params"`
}

type RPCResponse struct {
	ID      interface{} `json:"id"`
	JSONRPC string      `json:"jsonrpc"` // Should be "2.0"
	Error   *RPCError   `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

var didntMatch error = errors.New("Didn't match")

type cmd func(Args) (interface{}, error)

var stepRegister = []func(string, bool) (match bool, err error){}
var cmdRegister = map[string]cmd{}

// Execute JSON-RPC parameterized method-command,
// http://www.jsonrpc.org/specification
func ExecuteCMD(jsondata string) (RPCResponse, error) {
	var err error
	var ok bool
	var result interface{}
	var do cmd

	command := RPCMessage{}
	if err = json.Unmarshal([]byte(jsondata), &command); err != nil {
		err = errors.New("Can't parse message")
		return RPCResponse{nil, "2.0", &RPCError{ParseError, err.Error(), nil}, nil}, err
	}

	command.Params["gomateCMDId"] = strconv.Itoa(command.ID)

	if do, ok = cmdRegister[command.Method]; !ok {
		err = errors.New("Method not found")
		return RPCResponse{command.ID, "2.0", &RPCError{MethodNotFound, err.Error(), nil}, nil}, err
	}

	if result, err = do(Args(command.Params)); err != nil {
		return RPCResponse{command.ID, "2.0", &RPCError{InternalError, err.Error(), nil}, nil}, err
	}

	return RPCResponse{command.ID, "2.0", nil, result}, nil
}

// https://github.com/cucumber/cucumber/wiki/Given-When-Then
func stepImplementation(step string, do cmd, commands []string) {
	for _, command := range commands {
		cmdRegister[command] = do
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
