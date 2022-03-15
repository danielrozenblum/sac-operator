package converter

import (
	"testing"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"bitbucket.org/accezz-io/sac-operator/model"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConvertToModelWhenHTTP(t *testing.T) {
	// given
	converter := NewApplicationTypeConverter()

	applicationName := "test-application"
	siteID := "8cedca2d-3275-4f67-97be-088ece2a71b1"
	accessPolicies := []string{"8cedca2d-3275-4f67-97be-088ece2a71b1", "8cedca2d-3275-4f67-97be-088ece2a71b2"}

	application := buildTestApplication(applicationName, siteID, accessPolicies)

	// when
	result := converter.ConvertToModel(&application)

	// then
	assert.Nil(t, result.ID)
	assert.Equal(t, applicationName, result.Name)
	assert.Equal(t, model.DefaultType, result.Type)
	assert.Equal(t, model.DefaultSubType, result.SubType)
	assert.Equal(t, "http://net-tools:8080", result.InternalAddress)
	assert.Equal(t, siteID, result.SiteName)
	assert.Equal(t, accessPolicies, result.AccessPoliciesNames)
	assert.Equal(t, []string{}, result.ActivityPoliciesNames)
}

func TestConvertToModelWhenHTTPS(t *testing.T) {
	// given
	converter := NewApplicationTypeConverter()

	application := buildTestApplication("test-application", "12345", []string{})
	application.Spec.Service.Port = "443"

	// when
	result := converter.ConvertToModel(&application)

	// then
	assert.Nil(t, result.ID)
	assert.Equal(t, "test-application", result.Name)
	assert.Equal(t, model.DefaultType, result.Type)
	assert.Equal(t, model.DefaultSubType, result.SubType)
	assert.Equal(t, "https://net-tools:443", result.InternalAddress)
	assert.Equal(t, "12345", result.SiteName)
	assert.Equal(t, []string{}, result.AccessPoliciesNames)
	assert.Equal(t, []string{}, result.ActivityPoliciesNames)
}

func buildTestApplication(applicationName string, siteID string, accessPolicies []string) accessv1.Application {
	return accessv1.Application{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "Test-Namespace",
		},
		Spec: accessv1.ApplicationSpec{
			Type:    "",
			SubType: "",
			Service: accessv1.Service{
				Name: "net-tools",
				Port: "8080",
			},
			SiteName:              siteID,
			AccessPoliciesNames:   accessPolicies,
			ActivityPoliciesNames: []string{},
		},
		Status: accessv1.ApplicationStatus{Id: ""},
	}
}
