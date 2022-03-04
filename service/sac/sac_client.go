package sac

import (
	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
)

//go:generate mockery --name=SecureAccessCloudClient --inpackage --case=underscore --output=mockSecureAccessCloudClientInterface
type SecureAccessCloudClient interface {
	CreateApplication(applicationDTO *dto.ApplicationDTO) (*dto.ApplicationDTO, error)
	UpdateApplication(applicationDTO *dto.ApplicationDTO) (*dto.ApplicationDTO, error)
	FindApplicationByName(name string) (*dto.ApplicationDTO, error)
	FindApplicationByID(id string) (*dto.ApplicationDTO, error)
	DeleteApplication(id string) error

	FindPolicyByName(name string) (dto.PolicyDTO, error)
	FindPoliciesByNames(name []string) ([]dto.PolicyDTO, error)
	UpdatePolicies(applicationId string, applicationType model.ApplicationType, policies []string) error

	FindSiteByName(name string) (*dto.SiteDTO, error)
	CreateSite(siteDTO *dto.SiteDTO) (*dto.SiteDTO, error)
	DeleteSite(id string) error
	BindApplicationToSite(applicationId string, siteId string) error

	CreateConnector(siteDTO *dto.SiteDTO, connectorName string) (*dto.ConnectorObjects, error)
	ListConnectorsBySite(siteName string) ([]string, error)
	DeleteConnector(connectorID string) error
}
