package sac

import (
	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"github.com/google/uuid"
)

type SecureAccessCloudClient interface {
	CreateApplication(applicationDTO *dto.ApplicationDTO) (*dto.ApplicationDTO, error)
	UpdateApplication(applicationDTO *dto.ApplicationDTO) (*dto.ApplicationDTO, error)
	FindApplicationByName(name string) (*dto.ApplicationDTO, error)
	FindApplicationByID(id uuid.UUID) (*dto.ApplicationDTO, error)
	DeleteApplication(id uuid.UUID) error

	FindPolicyByName(name string) (dto.PolicyDTO, error)
	FindPoliciesByNames(name []string) ([]dto.PolicyDTO, error)
	UpdatePolicies(applicationId uuid.UUID, applicationType model.ApplicationType, policies []uuid.UUID) error

	FindSiteByName(name string) (*dto.SiteDTO, error)
	BindApplicationToSite(applicationId uuid.UUID, siteId uuid.UUID) error
}
