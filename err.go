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

// Errors type of error containing multiple entries for batch processing and collection of full set of errors
// (instead of failing on the first error)
type Errors interface {
	error
	Add(message string)
	AddError(e Error)
	Contains(e Error) bool
	Count() int
}

// Errors implementation
type errors struct {
	errors []Error
}

// NewErrors creates and returns a new Errors object
func NewErrors() Errors {
	return &errors{
		errors: []Error{},
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
func (es *errors) AddError(e Error) {
	es.errors = append(es.errors, e)
}

// Contains checks if the given error is part of the list of errors
func (es *errors) Contains(e Error) bool {
	for _, err := range es.errors {
		if err == e {
			return true
		}
	}
	return false
}

// Count gets the number of errors it contains
func (es *errors) Count() int {
	return len(es.errors)
}
