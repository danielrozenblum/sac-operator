package utils

import (
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"
)

func FromUIDType(id *types.UID) (*uuid.UUID, error) {
	if id == nil {
		return nil, nil
	}

	valueAsStr := string(*id)
	result, err := uuid.Parse(valueAsStr)
	return &result, err
}

func FromUUID(id uuid.UUID) *types.UID {
	var result types.UID
	result = types.UID(id.String())

	return &result
}

func FromString(valueAsStr string) (uuid.UUID, error) {
	return uuid.Parse(valueAsStr)
}

func FromInt64(i int64) *int64 {

	return &i

}

func ToStringArray(policyIDs []uuid.UUID) []string {
	var policyIDsAsString []string
	for _, policyID := range policyIDs {
		policyIDsAsString = append(policyIDsAsString, policyID.String())
	}

	return policyIDsAsString
}
