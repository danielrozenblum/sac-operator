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

func TestFromString(t *testing.T) {
	// given
	valueAsStr := "81399e0b-3981-4f3e-82cf-6c51067629fe"

	// when
	result, err := FromString(valueAsStr)

	// then
	assert.NoError(t, err)
	assert.Equal(t, valueAsStr, result.String())
}

func TestToStringArray(t *testing.T) {
	// given
	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	var uuidArray []uuid.UUID
	uuidArray = append(uuidArray, id1, id2, id3)

	// when
	result := ToStringArray(uuidArray)

	// then
	assert.Equal(t, 3, len(result))
	assert.Equal(t, id1.String(), result[0])
	assert.Equal(t, id2.String(), result[1])
	assert.Equal(t, id3.String(), result[2])
}

func TestToStringArrayWhenEmpty(t *testing.T) {
	// given
	var uuidArray []uuid.UUID

	// when
	result := ToStringArray(uuidArray)

	// then
	assert.Equal(t, 0, len(result))
}
