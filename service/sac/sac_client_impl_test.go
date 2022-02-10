package sac

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/util/rand"

	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

var sacClientTest *SecureAccessCloudClientTest

type SecureAccessCloudClientTest struct {
	client SecureAccessCloudClient
}

func (f *SecureAccessCloudClientTest) setup(t *testing.T) func(t *testing.T) {
	settings := &SecureAccessCloudSettings{
		//ClientID:     utils.GetMandatoryEnvironmentVariable(t, "SAC_CLIENT_ID"),
		//ClientSecret: utils.GetMandatoryEnvironmentVariable(t, "SAC_CLIENT_SECRET"),
		//TenantDomain: utils.GetMandatoryEnvironmentVariable(t, "SAC_TENANT_DOMAIN"),
		ClientID:     "a0090f01e9d27ffd82f4981c39a65392",
		ClientSecret: "f38b35a2e4eae78cbc85cc2ad745c6db5fdad6d179144279726ddfbd6112af02",
		TenantDomain: "symchatbotdemo.luminatesite.com",
	}

	sacClient := NewSecureAccessCloudClientImpl(settings)

	sacClientTest = &SecureAccessCloudClientTest{
		client: sacClient,
	}

	return func(t *testing.T) {
		// tearDown
	}
}

func TestFindApplicationByName(t *testing.T) {
	// given
	tearDown := sacClientTest.setup(t)
	defer tearDown(t)

	// when
	result, err := sacClientTest.client.FindApplicationByName("integration-test-application")

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, "integration-test-application", result.Name)
}

func TestFindApplicationByNameWhenNotFound(t *testing.T) {
	// given
	tearDown := sacClientTest.setup(t)
	defer tearDown(t)

	// when
	_, err := sacClientTest.client.FindApplicationByName("unknown-app")

	// then
	assert.Error(t, err)
	assert.Equal(t, ErrorNotFound, err)
}

func TestFindSiteByName(t *testing.T) {
	// given
	tearDown := sacClientTest.setup(t)
	defer tearDown(t)

	// when
	result, err := sacClientTest.client.FindSiteByName("integration-test-site")

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, "integration-test-site", result.Name)
}

func TestFindSiteByNameWhenNotFound(t *testing.T) {
	// given
	tearDown := sacClientTest.setup(t)
	defer tearDown(t)

	// when
	_, err := sacClientTest.client.FindSiteByName("unknown-site")

	// then
	assert.Error(t, err)
	assert.Equal(t, ErrorNotFound, err)
}

func TestCreateSite(t *testing.T) {
	// given
	tearDown := sacClientTest.setup(t)
	randomSiteName := fmt.Sprintf("create-site-%s", rand.String(4))
	site := &dto.SiteDTO{}
	defer func() {
		err := sacClientTest.client.DeleteSite(*site.ID)
		if err != nil {
			t.Errorf("failed deleteing site %+v", site)
		}
	}()
	defer tearDown(t)

	// when
	site, err := sacClientTest.client.CreateSite(&dto.SiteDTO{
		ID:   &uuid.UUID{},
		Name: randomSiteName,
	})

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, site)
	assert.Equal(t, randomSiteName, site.Name)

	// when
	connector := &dto.Connector{}
	connector, err = sacClientTest.client.CreateConnector(site)
	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, connector)
	assert.NotEmpty(t, connector.Otp)

}
