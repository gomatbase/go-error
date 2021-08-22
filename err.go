// Copyright 2021 GOM. All rights reserved.
// Since 25/06/2021 By GOM
// Licensed under MIT License

// Package err provides simple structures to handle error constants and error aggregations
package err

import (
	"bytes"
	"fmt"
)

// IError
// generic error interface adding a comparison function and parameters to be added to the error message when
// applicable.
type IError interface {
	// Inherits the error interface
	error
	// Equals
	// Compares the error with another. Must return true if the errors is functionally the same
	// (not necessarily the same message)
	Equals(error) bool
	// WithParameters
	// Takes parameters to build the error message. Only applicable for ErrorF implementations
	WithParameters(parameters ...interface{}) error
}

// Error Simple error type for log errors
type Error string

// Error returns the error message
func (e Error) Error() string {
	return string(e)
}

// Equals compares the error with another. Return true if it's the same message
func (e Error) Equals(err error) bool {
	return e == err
}

// WithParameters implemented because of IError. No actual use and returns itself
func (e Error) WithParameters(_ ...interface{}) error {
	return e
}

// errorF internal error type to keep errors with formatted messages and potential parameters.
// While it implements the error interface, so that it can be used to check if the error is part of an Errors collection,
// it is not meant to be used as a standard error, instead used as a pattern
// base for errors and can be used to compare if an errorF instance is considered to be the same error
type errorF struct {
	pattern string
}

// Error returns the error pattern
func (e *errorF) Error() string {
	return e.pattern
}

// Equals checks if the error is the same. Returns trus for the same object
func (e *errorF) Equals(err error) bool {
	return e == err
}

// WithParameters takes the parameters to create a new errorF instance which will have the parameters applied to the
// the errorF pattern, and will have the proper error message when printed.
func (e *errorF) WithParameters(parameters ...interface{}) error {
	return &errorFInstance{
		base: e,
		msg:  fmt.Sprintf(e.pattern, parameters...),
	}
}

// errorFInstance is the instance of an errorF containing the resolved error message for the specific error situation.
// This has the intention of providing error that may be compared for the type of error that it is (using the pattern constant
// while having situation specific messages
type errorFInstance struct {
	base *errorF
	msg  string
}

// Error returns the resolved pattern message with the given parameters from errorF
func (e *errorFInstance) Error() string {
	return e.msg
}

// Equals Checks if the error err is of the same type (same pattern) or if it's the same message
func (e *errorFInstance) Equals(err error) bool {
	return e.base == err || err.Error() == e.msg
}

// WithParameters shoud not be used, but for completeness it will provide anew errorFInstance with the same pattern
// and the message being the pattern resolved with the given parameters.
func (e *errorFInstance) WithParameters(parameters ...interface{}) error {
	return e.base.WithParameters(parameters...)
}

// ErrorF will return a new pattern based error which can be compared for type and produce situational error messages
func ErrorF(pattern string) IError {
	return &errorF{pattern: pattern}
}

// Errors type of error containing multiple entries for batch processing and collection of full set of errors
// (instead of failing on the first error)
type Errors interface {
	IError
	Add(message string)
	AddError(e IError)
	Contains(e IError) bool
	Count() int
}

// Errors implementation
type errors struct {
	errors []IError
}

// NewErrors creates and returns a new Errors object
func NewErrors() Errors {
	return &errors{
		errors: []IError{},
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

// Equals will check if the given error is the same object. Not meant to be used.
func (es *errors) Equals(e error) bool {
	// It's the same if it's the same collection of errors.
	return es == e
}

// WithParameters is a no-op and will return itself
func (es *errors) WithParameters(_ ...interface{}) error {
	return es
}

// Add add a new error message entry
func (es *errors) Add(message string) {
	es.AddError(Error(message))
}

// AddError add a new error entry
func (es *errors) AddError(e IError) {
	es.errors = append(es.errors, e)
}

// Contains checks if the given error is part of the list of errors
func (es *errors) Contains(e IError) bool {
	for _, err := range es.errors {
		if err.Equals(e) {
			return true
		}
	}
	return false
}

// Count gets the number of errors it contains
func (es *errors) Count() int {
	return len(es.errors)
}
