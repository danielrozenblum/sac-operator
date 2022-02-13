package sac

import (
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"github.com/google/uuid"
)

type SecureAccessCloudClient interface {
	CreateApplication(applicationDTO *dto.ApplicationDTO) (*dto.ApplicationDTO, error)
	UpdateApplication(applicationDTO *dto.ApplicationDTO) (*dto.ApplicationDTO, error)
	FindApplicationByName(name string) (*dto.ApplicationDTO, error)
	DeleteApplication(id uuid.UUID) error

	FindPolicyByName(name string) (dto.PolicyDTO, error)
	FindPoliciesByNames(name []string) ([]dto.PolicyDTO, error)
	UpdatePolicies(applicationId uuid.UUID, policies []uuid.UUID) error

	FindSiteByName(name string) (*dto.SiteDTO, error)
	CreateSite(siteDTO *dto.SiteDTO) (*dto.SiteDTO, error)
	DeleteSite(id uuid.UUID) error
	BindApplicationToSite(applicationId uuid.UUID, siteId uuid.UUID) error

	CreateConnector(siteDTO *dto.SiteDTO, connectorName string) (*dto.Connector, error)
	ListConnectorsBySite(siteName string) ([]dto.ConnectorPageDTO, error)
	DeleteConnector(connectorID uuid.UUID) error
}
