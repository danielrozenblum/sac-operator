package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetMandatoryEnvironmentVariable(t *testing.T, name string) string {
	value := os.Getenv(name)
	if value == "" {
		assert.FailNow(t, fmt.Sprintf("'%s' environment-variable not found. Have you forgot setting it?", name))
	}

	return value
}
