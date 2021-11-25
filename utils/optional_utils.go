package utils

func GetValueOrDefault(value interface{}, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}

	return value
}
