# Unbrokenwing

Use domain specific language (DSL) to verify your Go Lang Project.

For license information, see LICENSE.

## Introduction

Verifying software functionality and users scenarios are one of the
most critical steps in software engineering. Without the correct
behaviour, your project are lost, regardless how stable and bug free
it is.

Such verification needs to be automated, and continuously tested
by the environment. Therefor we need small and simple tools to
rapidly execute the test projects. Such tests are divided into
domain based scenarios and unit tests, Unbrokenwing provide
means to the former one.

### Example

A typical Hello World-example, providing users and project based
services might be structured like this:

```
.
├── Features
│   ├── projects
│   │   └── ...
│   │       ├── list.feature
│   │       └── remove.feature
│   └── users
│       ├── create.feature
│       ├── login.feature
│       └── step_definitions
│           └── create.go
├── hello.go
└── mypkg
    ├── world.go
    └── types.go
```

The Features directory contains one directory per feature area,
each area has then been divided into text files describing different
features to implement that area, such file needs to end with ".feature"
suffix in the name.

The .feature files contain descriptions written in Gherkin-ish
(https://github.com/cucumber/gherkin3) based language.

```
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

  Scenario: Change password for new user
    Given I'm successfully logged in as an admin in users pane
    And I change password for existing developer hacker from changeme to mysecretpwd
    When I press the update button
    Then only one user-record with name hacker should exist
    And user hacker should have password mysecretpwd

  Scenario: Remove a user
    Given I'm successfully logged in as an admin in users pane
    And I mark user hacker
    When I press the remove button
    Then no user named hacker should exist in user-record.
```

To test these scenarios we write a behaviour files which contains
calls to a limited set of Go Lang functions to test them.
Caller specifies a regular expression, and a callback functions as
first and second argument. Unbrokenwing makes use of the regular
expression to decide if a matching step in the scenario shall
trigger the callback. If callback returns nil, the test has succeeded.

create.go in our example above contains following lines:

```go
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
```

Lines beginning with package are irrelevant, and will be removed before
execution. All lines importing packages will be rearranged and placed
in the beginning of the executing code. Note that we return Pending
instances, which is a special type indicating that we have not finished
our test implementation yet.

To test the .feature file, you are able to run following
command:

```
unbrokenwing --pretty --dir ./create.feature test
```

In the terminal you will see a colorised version of the
test result. The output does also contain a status summerisation,
similar to:

```
3 scenario (0 undefined, 0 failures, 3 pending)
14 steps (10 undefined, 0 failures, 3 pending, 1 optout)
```

After that output you would see example code how to implement
Pending version of behaviours missing in todays setup.


## Requirements

- GoLang           (https://golang.org/)
- codegangsta/cli  (https://github.com/codegangsta/cli)

## Installation

```
go get github.com/dekelund/unbrokenwing
```

## Configuration

List available configurable values from CLI by running

```
unbrokenwing --help
```

Set default value by exporting environment variable in
the resource file belonging to your favorite shell e.g.,
$HOME/.bashrc

## Usage

```
NAME:
   unbrokenwing - Run behaviour driven tests as Gherik features

USAGE:
   unbrokenwing [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   feature-files			List feature files to STDOUT
   features					List features to STDOUT
   definitions, defs, code	List behaviours to STDOUT
   test, t					Tests either a test directory with features in it, or a .feature file
   help, h					Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --syslog						Redirect STDOUT to SysLog server
   --syslog-udp					Use UDP instead of TCP
   --syslog-raddr "localhost"	HOST/IP address to SysLog server
   --syslog-tag "unbrokenwing"	Tag output with specified text string
   --priority "6"				Log priority, use bitwised values from /usr/include/sys/syslog.h e.g.,
                                LOG_EMERG=0 LOG_ALERT=1 LOG_CRIT=2 LOG_ERR=3 LOG_WARNING=4 LOG_NOTICE=5 LOG_INFO=6 LOG_DEBUG=7
   --pretty						Print colorised result to STDOUT/STDERR
   --forensic					A kind of development mode, all generated files will be kept
   --dir "."					Relative path, to a feature-file or -directory (Current value: /Users/ekelund).

   --step-definitions "step_definitions"	Definitions folder name, should be located in features folder

   --help, -h					show help
   --version, -v				print the version
```

## Design decisions

Unbrokenwing has been divided into two subpackages:
unbrokenwing/compiler and unbrokenwing/unbrokenwing.

The former one contains code to parse features and behaviours.
This is used to generate executable code based on behaviour functions.

The later one contains help functions, required to execute tests
generated by unbrokenwing/compiler.

The main package, that binds together the two subpackages, parse
behaviour code and features. The result is either listed or executed
depending on subcommand provided in shell.

When subcommand "test" has been provided, it will go through following
steps:

1. Read behaviour code
2. Generate a new main-package/tool where setup function initialise callbacks and regular expression.
3. Write tool to disk
4. Compile tool
5. Execute tool, and pipe .feature file(s) to STDIN.

For educational and debugging purpose, you are able to print generated
code to STDOUT by running following command:

```
unbrokenwing --pretty code
```

## Known Problems

- Behaviour files with multiple packages in same import does not work

## Troubleshooting

## FAQ

-

## Maintainers

- Daniel Ekelund

## Missing features

* Data Tables
* Tags

## Alternatives

* http://goconvey.co
* https://onsi.github.io/ginkgo/
* https://github.com/cucumber/gherkin3

## Java Alternatives

* https://cucumber.io

## TODO

* Add Gherkin data tables support
* Add Gherkin tags support
* Verify SysLog implementation
