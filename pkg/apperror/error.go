// Package apperr implements the application errors.
package apperr

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
)

var (
	ErrLoginAlreadyExists              = errors.New("login already exists")
	ErrCardAlreadyExists               = errors.New("card with this number already exists")
	ErrUserNotFound                    = errors.New("user not found")
	ErrInvalidPassword                 = errors.New("invalid password")
	ErrUnableToGetUserLoginFromContext = errors.New("unable to get user login from context")
	ErrEmptyOrderRequest               = errors.New("empty order")
	ErrOrderUploadedByAnotherUser      = errors.New("order uploaded by another user")
	ErrOrderUploadedByUser             = errors.New("order uploaded by User")
	ErrOrderNotFound                   = errors.New("order not found")
	ErrBadNumber                       = errors.New("bad number")
	ErrNoOrders                        = errors.New("no orders")
	ErrRateLimit                       = errors.New("rate limit")
	ErrBalanceNotFound                 = errors.New("balance not found")
	ErrInsufficientFunds               = errors.New("insufficient funds")
	ErrNoWithdrawals                   = errors.New("no withdrawals")
)

// ValueError is an error that represents a value error.
type ValueError struct {
	caller  string
	message string
	err     error
}

// NewValueError creates a new ValueError with the given message, caller, and error.
//
// Parameters:
//
//	message string - the error message
//	caller string - the caller of the function tracing place in the code
//	err error - the original error
//
// Return type:
//
//	error - the newly created ValueError
func NewValueError(message string, caller string, err error) error {
	return &ValueError{
		caller:  caller,
		message: message,
		err:     err,
	}
}

// Error returns a string representing the error.
//
// No parameters.
// Returns a string.
func (v *ValueError) Error() string {
	return fmt.Sprintf("%s %s %s", v.caller, v.message, v.err)
}

// Unwrap returns the error that has been wrapped by ValueError.
// No parameters. Returns an error.
func (v *ValueError) Unwrap() error {
	return v.err
}

// Caller returns file name and line number of function call
func Caller() string {
	_, file, lineNo, ok := runtime.Caller(1)
	if !ok {
		return "runtime.Caller() failed"
	}

	fileName := path.Base(file)
	dir := filepath.Base(filepath.Dir(file))
	return fmt.Sprintf("%s/%s:%d", dir, fileName, lineNo)
}
