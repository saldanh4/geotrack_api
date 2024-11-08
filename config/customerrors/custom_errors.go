package customerrors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrDataBase       = errors.New("database error")
	ErrInternalServer = errors.New("internal server error")
)

type CustomError struct {
	BaseError error
	CustomMsg string
}

func (e *CustomError) Error() string {
	if e.CustomMsg != "" {
		return fmt.Sprintf("%s: %s", e.BaseError, e.CustomMsg)
	}
	return e.BaseError.Error()
}

func CustomErr(baseErr error, customMsg string) *CustomError {
	return &CustomError{
		BaseError: baseErr,
		CustomMsg: customMsg,
	}
}
