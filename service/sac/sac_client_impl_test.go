package sac

import (
	"fmt"
	"testing"

	"bitbucket.org/accezz-io/sac-operator/model"

	"bitbucket.org/accezz-io/sac-operator/utils"

	"k8s.io/apimachinery/pkg/util/rand"

	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"github.com/stretchr/testify/assert"
)

var sacClientTest *SecureAccessCloudClientTest

type SecureAccessCloudClientTest struct {
	client SecureAccessCloudClient
}

func (f *SecureAccessCloudClientTest) setup(t *testing.T) func(t *testing.T) {
	settings := &SecureAccessCloudSettings{
		ClientID:     utils.GetMandatoryEnvironmentVariable(t, "SAC_CLIENT_ID"),
		ClientSecret: utils.GetMandatoryEnvironmentVariable(t, "SAC_CLIENT_SECRET"),
		TenantDomain: utils.GetMandatoryEnvironmentVariable(t, "SAC_TENANT_DOMAIN"),
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
		err := sacClientTest.client.DeleteSite(site.ID)
		if err != nil {
			t.Errorf("failed deleteing site %+v", site)
		}
	}()
	defer tearDown(t)

	// when
	site, err := sacClientTest.client.CreateSite(&dto.SiteDTO{
		ID:   "",
		Name: randomSiteName,
	})

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, site)
	assert.Equal(t, randomSiteName, site.Name)

	// when
	connector := &dto.ConnectorObjects{}
	connector, err = sacClientTest.client.CreateConnector(site, "test")
	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, connector)
	assert.NotEmpty(t, connector.Otp)

}

func TestCreateApplication(t *testing.T) {
	// given
	tearDown := sacClientTest.setup(t)
	applicationName := fmt.Sprintf("create-application-%s", rand.String(4))
	application := &dto.ApplicationDTO{}
	defer func() {
		err := sacClientTest.client.DeleteApplication(application.ID)
		if err != nil {
			t.Errorf("failed deleteing application %+v", application)
		}
	}()
	defer tearDown(t)

	// when
	application, err := sacClientTest.client.CreateApplication(&dto.ApplicationDTO{
		Name: applicationName,
		Type: model.HTTP,
		ConnectionSettings: dto.ConnectionSettingsDTO{
			InternalAddress: "https://httpbin.org/",
		},
		Icon:                  "",
		IsVisible:             false,
		IsNotificationEnabled: false,
		Enabled:               false,
	})

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, application)
	assert.Equal(t, applicationName, application.Name)

}

func TestBindApplicationToSite(t *testing.T) {
	// given
	tearDown := sacClientTest.setup(t)
	applicationName := fmt.Sprintf("operator-application-%s", rand.String(4))
	siteName := fmt.Sprintf("operator-site-%s", rand.String(4))
	application := &dto.ApplicationDTO{}
	site := &dto.SiteDTO{}
	defer func() {
		err := sacClientTest.client.DeleteSite(site.ID)
		if err != nil {
			t.Errorf("failed deleteing site %+v", application)
		}
		err = sacClientTest.client.DeleteApplication(application.ID)
		if err != nil {
			t.Errorf("failed deleteing application %+v", application)
		}
	}()
	defer tearDown(t)

	// when
	application, err := sacClientTest.client.CreateApplication(&dto.ApplicationDTO{
		Name: applicationName,
		Type: model.HTTP,
		ConnectionSettings: dto.ConnectionSettingsDTO{
			InternalAddress: "https://httpbin.org/",
		},
	})
	assert.NoError(t, err)
	site, err = sacClientTest.client.CreateSite(&dto.SiteDTO{
		Name: siteName,
	})
	assert.NoError(t, err)
	err = sacClientTest.client.BindApplicationToSite(application.ID, site.ID)
	assert.NoError(t, err)
	// then

}
