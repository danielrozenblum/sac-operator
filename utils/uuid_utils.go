package utils

import (
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"
)

func FromUIDType(id *types.UID) (uuid.UUID, error) {
	valueAsStr := string(*id)

	return uuid.Parse(valueAsStr)
}

func FromUUID(id uuid.UUID) *types.UID {
	var result types.UID
	result = types.UID(id.String())

	return &result
}
