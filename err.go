// Copyright 2021 GOM. All rights reserved.
// Since 25/06/2021 By GOM
// Licensed under MIT License

// Package err provides simple structures to handle error constants and error aggregations
package err

import (
	"bytes"
	"fmt"
)

// Error Simple error type for log errors
type Error string

// Error returns the error message
func (e Error) Error() string {
	return string(e)
}

// ErrorF an error Kind
type ErrorF string

// Error returns the error message
func (ef ErrorF) Error() string {
	return string(ef)
}

// IsKindOf checks if the err is a kind of then ErrorF
func (ef ErrorF) IsKindOf(e error) bool {
	if efi, isErrorFInstance := e.(*errorFInstance); isErrorFInstance {
		return efi.base == ef
	}
	return false
}

func (ef ErrorF) WithValues(values ...interface{}) *errorFInstance {
	return &errorFInstance{
		base: ef,
		msg:  fmt.Sprintf(string(ef), values...),
	}
}

// errorFInstance is the instance of an ErrorF containing the resolved error message for the specific error situation.
// This has the intention of providing errors that may be compared for the kind of error that it is (using the pattern constant)
// while having situation specific messages
type errorFInstance struct {
	base ErrorF
	msg  string
}

// Error returns the resolved pattern message with the given parameters from errorF
func (e *errorFInstance) Error() string {
	return e.msg
}

// Errors type of error containing multiple entries for batch processing and collection of full set of errors
// (instead of failing on the first error)
type Errors interface {
	error
	Add(message string)
	AddError(e error)
	Contains(e error) bool
	Count() int
}

// Errors implementation
type errors struct {
	errors []error
}

// NewErrors creates and returns a new Errors object
func NewErrors() Errors {
	return &errors{
		errors: []error{},
	}
}

// Error returns the error message by joining all the existing errors and separating them by new lines
func (es *errors) Error() string {
	buffer := &bytes.Buffer{}
	for _, message := range es.errors {
		// explicitly ignore the printing errors as none is expected
		_, _ = fmt.Fprintln(buffer, message)
	}
	return buffer.String()
}

// Add add a new error message entry
func (es *errors) Add(message string) {
	es.AddError(Error(message))
}

// AddError add a new error entry
func (es *errors) AddError(e error) {
	es.errors = append(es.errors, e)
}

// Contains checks if the given error is part of the list of errors
func (es *errors) Contains(e error) bool {
	isError := func(err error) bool {
		return e == err
	}
	if errorKind, isKind := e.(ErrorF); isKind {
		isError = func(err error) bool {
			return errorKind.IsKindOf(err)
		}
	}
	for _, err := range es.errors {
		if isError(err) {
			return true
		}
	}
	return false
}

// Count gets the number of errors it contains
func (es *errors) Count() int {
	return len(es.errors)
}
