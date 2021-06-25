# go-error

err is a package with utility types to handle errors as constants and collections.

The idea is to define errors as string constants, using simple comparisons to check the kind
of error raised and having an error aggregator when an operation may return more than one error
(for example in validation functions) but still allowing individual errors that may have occurred.

A simplistic usage is shown in the following example

## Example

```go
package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gomatbase/go-error"
)

var (
	greetings = make(map[string]string)
	people    = make(map[string]bool)
)

const (
	errIncorrectTime = err.Error("incorrect time of day")
	errInvalidName   = err.Error("invalid name format")
	errUnknownPerson = err.Error("unknown person")
)

func init() {
	greetings["morning"] = "Good Morning"
	greetings["afternoon"] = "Good Afternoon"
	greetings["evening"] = "Good Evening"
	people["Jack"] = true
	people["Mary"] = true
}

func getGreeting(time string) (string, error) {
	if greeting, found := greetings[time]; found {
		return greeting, nil
	}
	return "", errIncorrectTime
}

func validateName(name string) error {
	e := err.NewErrors()
	r, _ := regexp.Compile("[A-Z][a-z]*")
	if !r.MatchString(name) {
		e.AddError(errInvalidName)
	}
	if _, found := people[name]; !found {
		e.AddError(errUnknownPerson)
	}
	if e.Count() > 0 {
		return e
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Invalid Number of Parameters")
		os.Exit(1)
	}
	greeting, e := getGreeting(os.Args[1])
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	name := os.Args[2]
	if e := validateName(name); e != nil {
		errors := e.(err.Errors)
		if errors.Contains(errUnknownPerson) && errors.Count() == 1 {
			name = "Unknown Person"
		} else {
			fmt.Print(e)
			os.Exit(1)
		}
	}

	fmt.Println(greeting, name)
}
```