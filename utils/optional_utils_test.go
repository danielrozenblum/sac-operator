package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetValueOrDefault(t *testing.T) {
	// given
	value := "value"
	defaultValue := "default"

	// when
	result := GetValueOrDefault(value, defaultValue)

	// then
	assert.NotEmpty(t, result)
	assert.Equal(t, value, result)
}

func TestGetValueOrDefaultWhenNil(t *testing.T) {
	// given
	defaultValue := "default"

	// when
	result := GetValueOrDefault(nil, defaultValue)

	// then
	assert.NotEmpty(t, result)
	assert.Equal(t, defaultValue, result)
}

func TestGetValueOrDefaultWhenDefaultNil(t *testing.T) {
	// given

	// when
	result := GetValueOrDefault(nil, nil)

	// then
	assert.Nil(t, result)
}
