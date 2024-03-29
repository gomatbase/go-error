// Copyright 2021 GOM. All rights reserved.
// Since 25/06/2021 By GOM
// Licensed under MIT License

package err_test

import (
	"fmt"
	"testing"

	"github.com/gomatbase/go-error"
)

const (
	sampleError1 = err.Error("sample error 1")
	sampleError2 = err.Error("sample error 2")
	sampleError3 = err.Error("sample error 3")
	sampleErrorF = err.ErrorF("test %s")
)

func TestError(t *testing.T) {
	e := err.Error("test")
	if e.Error() != "test" {
		t.Error("Error message not expected", e.Error())
	}
}

func TestErrorF(t *testing.T) {
	if sampleErrorF.Error() != "test %s" {
		t.Error("Error message not expected", sampleErrorF.Error())
	}

	e := sampleErrorF.WithValues("something")
	if e.Error() != "test something" {
		t.Error("Error message not expected", e.Error())
	}
	if !sampleErrorF.IsKindOf(e) {
		t.Error("Error is not of the expected kind")
	}
	if !err.ErrorF("test %s").IsKindOf(e) {
		t.Error("Error is not of the expected kind")
	}
}

func TestErrors(t *testing.T) {
	es := err.Errors()
	es.AddError(sampleError1)
	es.Add("sample error 2")

	es.AddError(sampleErrorF.WithValues("something"))

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
		if !es.Contains(err.Error("sample error 1")) {
			t.Error("Errors doesn't contain sample 1")
		}
		if !es.Contains(err.Error("sample error 2")) {
			t.Error("Errors doesn't contain sample 2")
		}
		if es.Contains(err.Error("sample error 3")) {
			t.Error("Errors contains sample 3")
		}
		if !es.Contains(err.ErrorF("test %s")) {
			t.Error("Errors doesn't containd errorf")
		}
		if es.Contains(err.Error("test something")) {
			t.Error("Errors identifies errorF instance as an error instance")
		}
	})
}

func TestIsContainedIn(t *testing.T) {
	es := err.Errors()
	es.AddError(sampleError1)
	efInstance := sampleErrorF.WithValues("something")
	es.AddError(efInstance)

	if !err.IsContainedIn(sampleError1, es) {
		t.Error("sample error 1 should be contained in errors")
	}
	if err.IsContainedIn(sampleError2, es) {
		t.Error("sample error 2 should not be contained in errors")
	}
	if !err.IsContainedIn(sampleErrorF, es) {
		t.Error("sample error f should be contained in errors")
	}
	if !err.IsContainedIn(sampleErrorF, efInstance) {
		t.Error("sample error f should be contained in error f instance")
	}
	if !err.IsContainedIn(sampleError1, err.Error("sample error 1")) {
		t.Error("sample error 1 should be contained in sample error 1")
	}
	if err.IsContainedIn(sampleError1, sampleError2) {
		t.Error("sample error 1 should not be contained in sample error 2")
	}
}

func TestCount(t *testing.T) {
	es := err.Errors()
	es.AddError(sampleError1)
	es.AddError(sampleError2)

	if err.Count(nil) != 0 {
		t.Error("null errors should result in a 0 count")
	}
	if count := err.Count(es); count != 2 {
		t.Error("Count of errors not the expected amount :", count)
	}
	if err.Count(sampleError2) != 1 {
		t.Error("Count of any error should result in 1")
	}
}

func ExampleErrors() {
	es := err.Errors()
	es.AddError(sampleError1)
	es.AddError(sampleError2)
	es.AddError(sampleError3)
	es.AddError(sampleErrorF.WithValues("something"))

	fmt.Print(es.Error())

	// Output:
	//
	// sample error 1
	// sample error 2
	// sample error 3
	// test something
}
