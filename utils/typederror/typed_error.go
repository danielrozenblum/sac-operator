package typederror

import "github.com/pkg/errors"

type ErrorType error

var (
	UnknownError          = errors.New("UnknownError")
	UnrecoverableError    = errors.New("UnrecoverableError")
	PartiallySuccessError = errors.New("PartiallySuccessError")
)
