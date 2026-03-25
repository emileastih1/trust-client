package errors

import (
	"github.com/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// Internal Error.
	Internal Code = 0

	// The passed argument does not pass validation checks.
	InvalidArgument Code = 1

	// The resource already exists.
	AlreadyExists Code = 2

	// The resource could not be found.
	NotFound Code = 3

	// The current user does not have permissions to perform the request.
	PermissionDenied Code = 4

	// Authentication failed for the user.
	Unauthenticated Code = 5

	// The functionality has not yet been implemented.
	Unimplemented Code = 6
)

const ValidationErrorPrefix = "validation failure:"

type Code uint32

type Error struct {
	Code Code
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func EqualErrorCode(err error, code Code) bool {
	//nolint:errorlint
	re, ok := err.(*Error)
	if !ok {
		return false
	}
	return re.Code == code
}

func New(message string, code Code) *Error {
	underlyingError := errors.New(message)
	newError := &Error{
		Code: code,
		Err:  underlyingError,
	}
	return newError
}

// AsError returns the first error in err's chain that matches the internal Error type, and returns nil if not found.
func AsError(err error) *Error {
	var e *Error
	errors.As(err, &e)
	return e
}

func Wrap(err error, message string, code Code) error {
	underlyingError := errors.Wrap(err, message)
	newError := &Error{
		Code: code,
		Err:  underlyingError,
	}
	return newError
}

// Maps service errors to gRPC errors to be used by the request interceptor.
func ToGrpcError(err error) error {
	if err == nil {
		return nil
	}

	// Default to using Internal error
	code := codes.Internal

	if e := AsError(err); e != nil {
		switch e.Code {
		case Internal:
			code = codes.Internal
		case InvalidArgument:
			code = codes.InvalidArgument
		case AlreadyExists:
			code = codes.AlreadyExists
		case NotFound:
			code = codes.NotFound
		case PermissionDenied:
			code = codes.PermissionDenied
		case Unauthenticated:
			code = codes.Unauthenticated
		case Unimplemented:
			code = codes.Unimplemented
		}
	}
	// this is to hide the internal error message so that we do not leak implementation details to callers.
	if code == codes.Internal {
		return status.Error(code, "internal service error.")
	}
	if code == codes.InvalidArgument {
		return status.Error(code, ValidationErrorPrefix+" "+err.Error())
	}
	return status.Error(code, err.Error())
}
