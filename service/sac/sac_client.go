package sac

import "bitbucket.org/accezz-io/sac-operator/service/sac/dto"

type SecureAccessCloudClient interface {
	CreateApplication() error
	UpdateApplication() error
	FindApplicationByName(name string) (dto.Application, error)
	DeleteApplication(id string) error

	FindPolicyByName(name string) (dto.Policy, error)
	AddApplicationToPolicy() error
	RemoveApplicationFromPolicy() error

	FindSiteByName(name string) (dto.Site, error)
	AddApplicationToSite() error
	RemoveApplicationFromSite() error
}
