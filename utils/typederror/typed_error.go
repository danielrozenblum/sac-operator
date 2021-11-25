package typederror

import "fmt"

type ErrorType string

const (
	UnknownError          ErrorType = "UnknownError"
	UnrecoverableError              = "UnrecoverableError"
	PartiallySuccessError           = "PartiallySuccessError"
)

func (e ErrorType) String() string {
	return e.String()
}

type TypedError struct {
	Type ErrorType
	Err  error
}

func WrapError(errorType ErrorType, err error) error {
	return TypedError{Type: errorType, Err: err}
}

func (e TypedError) Unwrap() error { return e.Err }

func (e TypedError) Error() string {
	return fmt.Sprintf("error %v (type: '%s')", e.Err, e.Type)
}

func IsErrorType(errorType ErrorType, err error) bool {
	typedError, ok := err.(TypedError)
	if !ok {
		return false
	}

	return typedError.Type == errorType
}
