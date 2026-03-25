package errors_test

import (
	"errors"
	"fmt"
	"testing"

	serviceErrors "bulletin-board-api/internal/errors"

	"github.com/stretchr/testify/assert"
)

var (
	errNew = errors.New("errors.New() non-service error")
	errFmt = fmt.Errorf("fmt.Errorf() non-service error")
)

func TestErrors_AsServiceError(t *testing.T) {
	tests := []struct {
		err  error
		code serviceErrors.Code
		msg  string
	}{
		{
			err:  serviceErrors.New("test error", serviceErrors.NotFound),
			code: serviceErrors.NotFound,
			msg:  "test error",
		},
		{
			err:  serviceErrors.New("test error 2", serviceErrors.Unimplemented),
			code: serviceErrors.Unimplemented,
			msg:  "test error 2",
		},
		{
			err:  serviceErrors.Wrap(errFmt, "test error 3", serviceErrors.Unauthenticated),
			code: serviceErrors.Unauthenticated,
			msg:  "test error 3: fmt.Errorf() non-service error",
		},
		{
			err:  serviceErrors.Wrap(errNew, "test error 4", serviceErrors.PermissionDenied),
			code: serviceErrors.PermissionDenied,
			msg:  "test error 4: errors.New() non-service error",
		},
	}

	for _, test := range tests {
		got := serviceErrors.AsError(test.err)
		assert.NotNil(t, got)
		if !assert.EqualError(t, got, test.msg) {
			t.Errorf("unexpected error, got: %v, want %v", got, test.msg)
		}

		if !serviceErrors.EqualErrorCode(got, test.code) {
			t.Errorf("unexpected error code, got: %v, want %v", got.Code, test.code)
		}
	}
}

func TestErrors_AsNonServiceError(t *testing.T) {
	tests := []struct {
		err error
	}{
		{err: errNew},
		{err: errFmt},
		{err: nil},
	}

	for _, test := range tests {
		got := serviceErrors.AsError(test.err)
		assert.Nil(t, got)
	}
}
