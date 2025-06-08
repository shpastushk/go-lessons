package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	result := fmt.Sprintf("%d errors occured:\n\t", len(e.errors))
	values := make([]string, 0)
	for _, err := range e.errors {
		values = append(values, "* "+err.Error())
	}
	if len(values) != 0 {
		result += strings.Join(values, "\t") + "\n"
	}
	return result
}

func Append(err error, errs ...error) *MultiError {
	switch err := err.(type) {
	case *MultiError:
		if err == nil {
			err = &MultiError{}
		}
		err.errors = append(err.errors, errs...)
		return err
	default:
		mErr := &MultiError{}
		mErr.errors = append(mErr.errors, errs...)
		return mErr
	}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
