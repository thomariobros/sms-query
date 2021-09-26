package errors

import (
	"fmt"
)

const (
	NOT_ALLOWED = 1
	WRONG_INPUT = 2
)

type CustomError struct {
	Code int
}

func NewCustomError(code int) CustomError {
	return CustomError{
		Code: code,
	}
}

func (e CustomError) Error() string {
	return fmt.Sprintf("%d", e.Code)
}
