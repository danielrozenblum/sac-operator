package typederror

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrapPartiallySuccessError(t *testing.T) {
	// given
	innerErr := errors.New("unit-test")

	// when
	wrappedError := WrapError(PartiallySuccessError, innerErr)

	// then
	assert.True(t, IsErrorType(PartiallySuccessError, wrappedError))
}

func TestIsErrorTypeWhenNotSameType(t *testing.T) {
	// given
	innerErr := errors.New("unit-test")

	// when
	wrappedError := WrapError(PartiallySuccessError, innerErr)

	// then
	assert.False(t, IsErrorType(UnrecoverableError, wrappedError))
}

func TestIsErrorTypeWhenNotTypedError(t *testing.T) {
	// given
	innerErr := errors.New("unit-test")

	// when

	// then
	assert.False(t, IsErrorType(UnrecoverableError, innerErr))
}
