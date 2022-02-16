package utils

import "bitbucket.org/accezz-io/sac-operator/model"

func GetValueOrDefault(value interface{}, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}

	return value
}

func GetStringPtrValueOrDefault(value *string, defaultValue string) string {
	if value == nil {
		return defaultValue
	}

	return *value
}

func GetApplicationTypeOrDefault(value *model.ApplicationType, defaultValue model.ApplicationType) model.ApplicationType {
	if value == nil {
		return defaultValue
	}

	return *value
}

func GetApplicationSubTypeOrDefault(value *model.ApplicationSubType, defaultValue model.ApplicationSubType) model.ApplicationSubType {
	if value == nil {
		return defaultValue
	}

	return *value
}
