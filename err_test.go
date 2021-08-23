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

var (
	sampleErrorF = ErrorF("test %s")
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

func TestErrorF(t *testing.T) {
	ef := ErrorF("test %s")
	ef2 := ErrorF("test %s")
	if ef.Equals(ef2) {
		t.Error("Error kinds should only be equal to itself")
	}
	if !ef.Equals(ef) {
		t.Error("Error kinds should be equal to itself")
	}
	e := ef.WithParameters("something").(IError)
	if e.Error() != "test something" {
		t.Error("Error message not expected", e.Error())
	}
	if !e.Equals(Error("test something")) {
		t.Error("Error equals should succeed on errors with the same message")
	}
	if !e.Equals(ef) {
		t.Error("Error equals should succeed comparing with the kind of error")
	}
	if e.Equals(Error("test other")) {
		t.Error("Error equals should fail for any other error")
	}
	if e.Equals(ef2) {
		t.Error("Error equals should fail comparing with any other error kind")
	}
	if e.WithParameters("something") == e {
		t.Error("With parameters should result in a different error instance")
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

	es.AddError(sampleErrorF.WithParameters("something").(IError))

	if es.Count() != 3 {
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
		if !es.Contains(sampleErrorF) {
			t.Error("Errors doesn't contain sample errorF")
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
		if !es.Contains(Error("test something")) {
			t.Error("Errors doesn't containt errorf")
		}
	})
}

func ExampleErrors() {
	es := NewErrors()
	es.AddError(sampleError1)
	es.AddError(sampleError2)
	es.AddError(sampleError3)
	es.AddError(sampleErrorF.WithParameters("something").(IError))

	fmt.Print(es.Error())

	// Output:
	//
	// sample error 1
	// sample error 2
	// sample error 3
	// test something
}
