// Copyright 2021 GOM. All rights reserved.
// Since 25/06/2021 By GOM
// Licensed under MIT License

package err

import (
	"fmt"
	"testing"
)

const (
	sampleError1 = Error("sample error 1")
	sampleError2 = Error("sample error 2")
	sampleError3 = Error("sample error 3")
)

func TestError(t *testing.T) {
	e := Error("test")
	if e.Error() != "test" {
		t.Error("Error message not expected", e.Error())
	}
	if !e.Equals(Error("test")) {
		t.Error("Error equals should succeed on errors with the same message")
	}
	if e.Equals(Error("test other")) {
		t.Error("Error equals should fail for any other error")
	}
	if e.WithParameters("something") != e {
		t.Error("Plain Error should have WithParameters as a neutral operation")
	}
}

func TestErrors(t *testing.T) {
	es := NewErrors()
	es.AddError(sampleError1)
	es.Add("sample error 2")
	if !es.Equals(es) {
		t.Error("Errors equals should succeed with itself")
	}
	es2 := NewErrors()
	es2.AddError(sampleError1)
	es2.Add("sample error 2")
	if es.Equals(es2) {
		t.Error("Errors equals should only succeed with itself")
	}
	if es.WithParameters("something") != es {
		t.Error("Errors should have WithParameters as a neutral operation")
	}
	if es.Count() != 2 {
		t.Error("Reporting incorrect number of errors :", es.Count())
	}
	t.Run("Test if errors contains an error", func(t *testing.T) {
		if !es.Contains(sampleError1) {
			t.Error("Errors doesn't contain sample 1")
		}
		if !es.Contains(sampleError2) {
			t.Error("Errors doesn't contain sample 2")
		}
		if es.Contains(sampleError3) {
			t.Error("Errors contains sample 3")
		}
	})

	t.Run("Test if errors contains an error message", func(t *testing.T) {
		if !es.Contains(Error("sample error 1")) {
			t.Error("Errors doesn't contain sample 1")
		}
		if !es.Contains(Error("sample error 2")) {
			t.Error("Errors doesn't contain sample 2")
		}
		if es.Contains(Error("sample error 3")) {
			t.Error("Errors contains sample 3")
		}
	})
}

func ExampleErrors() {
	es := NewErrors()
	es.AddError(sampleError1)
	es.AddError(sampleError2)
	es.AddError(sampleError3)

	fmt.Print(es.Error())

	// Output:
	//
	// sample error 1
	// sample error 2
	// sample error 3
}
