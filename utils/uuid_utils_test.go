package utils

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

func TestFromValidUIDType(t *testing.T) {
	// given
	idStr := "5b46eda0-4a7c-4ba0-8376-a842b03a2165"
	uid := types.UID(idStr)

	// when
	result, err := FromUIDType(&uid)

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, idStr, result.String())
}

func TestFromInvalidUIDType(t *testing.T) {
	// given
	uid := types.UID("invalid")

	// when
	_, err := FromUIDType(&uid)

	// then
	assert.Error(t, err)
}

func TestFromUUID(t *testing.T) {
	// given
	id := uuid.New()

	// when
	result := FromUUID(id)

	// then
	assert.NotEmpty(t, result)
	assert.Equal(t, id.String(), string(*result))
}
