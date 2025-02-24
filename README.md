> Project is archived in deference to https://github.com/gomatbase/csn

# go-error

err is a package with utility types to handle errors as constants and collections.

The idea is to define errors as string constants, using simple comparisons to check the kind
of error raised and having an error aggregator when an operation may return more than one error
(for example in validation functions) but still allowing individual errors that may have occurred.

## Error

Defining simple static errors is as simple as defining a string constant

```go
const errIncorrectTime = err.Error("incorrect time of day")
```

As an error it can be simply returned.

## ErrorF

ErrorF introduces the concept of "Kind of Error". This basically serves to define a base type of error with custom messages related to the
context where the error was raised, while still being able to identify that an error of a certain kind has been returned. The use is mainly
for logging errors or using the error message for some operation, while still being able to handle the error at a generic level.

ErrorF itself is only meant to be used to generate error instances with custom messages (See Caveats).

```go
const errUnknownPerson = err.ErrorF("Unknown Person (%s)")

func doSomething() error {
	return errUnknownPerson.WithValues("John Doe")
}

func main() {
	if e := doSomething; errUnknownPerson.IsKindOf(e) {
	    fmt.Println(e.Error()) // will print Unknown Person (John Dow)
    }
}
```

## Errors

Errors is a collection of errors. It is meant to aggregate more than one error that a function may return when each one is not meant to
be critical but the total amount of errors might. This would be the example for a validation function where it is desired to test all
the validation failures instead of returning error on the first error occurrence.

Checking for error containment works both for Error and ErrorF kind of errors.

```go
const (
    errIncorrectTime = err.Error("incorrect time of day")
    errUnknownPerson = err.ErrorF("Unknown Person (%s)")
)

func doSomething() error {
	e := err.Errors()
	e.AddError(errIncorrectTime)
	e.AddError(errUnknownPerson.WithValues("Someone"))
	e.Add("some error")
	return e
}

func main() {
	if e := doSomething; e != nil {
		errors := e.(err.IErrors)
		fmt.Println(errors.Contains(errIncorrectTime)) // prints true 
		fmt.Println(errors.Contains(errUnknownPerson)) // prints true 
		fmt.Println(errors.Contains(err.Error("some error"))) // prints true 
		fmt.Println(errors.Contains(err.Error("some other error"))) // prints false
	}
}
```

When checking for the presence of an error in a collection it is also possible to use the utility function IsContainedIn,
where the main function from the previous example would be. The utility function Count(e error) int can also be used to 
count errors of a potential Errors object without the need to assert its type first.

```go
func main() {
	if e := doSomething; e != nil {
		errors := e.(err.IErrors)
		fmt.Println(IsContainedIn(errIncorrectTime,errors)) // prints true 
		fmt.Println(IsContainedIn(errUnknownPerson,errors)) // prints true 
		fmt.Println(IsContainedIn(err.Error("some error"),errors)) // prints true 
		fmt.Println(IsContainedIn(err.Error("some other error"),errors)) // prints false
	}
}
```


## Example
A simplistic usage is shown in the following example

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
	errUnknownPerson = err.ErrorF("Unknown Person (%s)")
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
	e := err.Errors()
	r, _ := regexp.Compile("[A-Z][a-z]*")
	if !r.MatchString(name) {
		e.AddError(errInvalidName)
	}
	if _, found := people[name]; !found {
		e.AddError(errUnknownPerson.WithValues(name))
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
		if err.IsContainedIn(errUnknownPerson,e) && err.Count(e) == 1 {
			name = e.Error()
		} else {
			fmt.Print(e)
			os.Exit(1)
		}
	}

	fmt.Println(greeting, name)
}
```

## Caveats
Some details to keep in mind

### ErrorF as an error
ErrorF implements the error interface. This means it can be used as a standard error which will result in having an error message that
will be the string format pattern. This also means that there is nothing preventing ErrorF to be used as a standard error constant, except
that trying to find if that specific error is contained in an Errors collection will always fail as ErrorF will always try to compare
to errors of which it is a kind of (will only match errorFInstances).

### ErrorF (Non-)Uniqueness
ErrorF is in fact a string that represents a format pattern. This means that the rules for string comparison apply to it. This also means
that any identical ErrorFs, even though they might be defined as two distinct constants (or variables) they are considered to be the same,
and when checking for errors of the same kind, any errorF instance having the same ErrorF base format pattern will be considered to be of
the same kind (check err_test for an example where that property is tested).

### ErrorF instances.
ErrorF instances are objects referencing the originating string format pattern, and the actual values that should be applied to the pattern.
The number of values in the pattern are of the responsibility of the developer. Instantiating a new ErrorF with not enough or too many values
will not raise any error but will result into an error with the corresponding "(%!s(MISSING))" tokens for not enough values and
"%!(EXTRA string=)" tokens for too many.

### ErrorF containment checks for Errors
When checking the existence of an error within an Errors collection, if using an actual ErrorFInstance to check, it will always fail unless
it's an identical errorFInstance (either itself or an errorFInstance with the same string format pattern and the same actual error message).
This works as expected and errorFInstances are usually meant to be tested for their kindness.
