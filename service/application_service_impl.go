package service

import (
	"context"
	"errors"
	"fmt"

	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service/sac"
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"bitbucket.org/accezz-io/sac-operator/utils/typederror"
	"github.com/google/uuid"
	logger "sigs.k8s.io/controller-runtime/pkg/log"
)

type ApplicationServiceImpl struct {
	sacClient sac.SecureAccessCloudClient
}

func NewApplicationServiceImpl(sacClient sac.SecureAccessCloudClient) ApplicationService {
	return &ApplicationServiceImpl{sacClient: sacClient}
}

func (a *ApplicationServiceImpl) Create(ctx context.Context, applicationToCreate *model.Application) (*model.Application, error) {
	log := logger.FromContext(ctx)
	log.Info("Trying to create application: " + applicationToCreate.String())

	// 1. Find Application by Name to verify the name isn't used
	_, err := a.sacClient.FindApplicationByName(applicationToCreate.Name)
	if err != nil && err != sac.ErrorNotFound {
		return &model.Application{}, err
	}

	if err == nil {
		return &model.Application{},
			typederror.WrapError(
				typederror.UnrecoverableError,
				errors.New(fmt.Sprintf("application with name: '%s' already exists.", applicationToCreate.Name)),
			)
	}

	// 2. Validate Site Exists
	site, err := a.sacClient.FindSiteByName(applicationToCreate.Site)
	if err != nil {
		return nil, typederror.WrapError(typederror.UnrecoverableError, err)
	}

	// 3. Validate Policies Exists
	policies, err := a.sacClient.FindPoliciesByNames(append(applicationToCreate.AccessPolicies, applicationToCreate.ActivityPolicies...))
	if err != nil {
		return nil, typederror.WrapError(typederror.UnrecoverableError, err)
	}

	// 4. createSite Application
	applicationDTO := dto.FromApplicationModel(applicationToCreate)
	createdApplicationDTO, err := a.sacClient.CreateApplication(applicationDTO)
	if err != nil {
		return nil, err
	}

	// 5. Convert back to model
	result := dto.ToApplicationModel(createdApplicationDTO, site.ID)

	// 6. Bind Site & Policies (Idempotent)
	err = a.bindSiteAndPolicies(result, site, policies)
	if err != nil {
		return nil, err
	}

	log.Info("Application: '" + applicationToCreate.String() + "' created successfully.")
	return result, nil
}

func (a *ApplicationServiceImpl) Update(ctx context.Context, updatedApplication *model.Application) (*model.Application, error) {
	log := logger.FromContext(ctx)
	log.Info("Dummy Update Application")

	// 1. Validate Site Exists
	site, err := a.sacClient.FindSiteByName(updatedApplication.Site)
	if err != nil {
		return nil, typederror.WrapError(typederror.UnrecoverableError, err)
	}

	// 2. Validate Policies Exists
	policies, err := a.sacClient.FindPoliciesByNames(append(updatedApplication.AccessPolicies, updatedApplication.ActivityPolicies...))
	if err != nil {
		return nil, typederror.WrapError(typederror.UnrecoverableError, err)
	}

	// 3. Update Application
	applicationDTO, err := a.CompleteApplication(ctx, updatedApplication)
	if err != nil {
		return nil, err
	}

	updatedApplicationDTO, err := a.sacClient.UpdateApplication(applicationDTO)
	if err != nil {
		return nil, err
	}

	// 4. Convert back to model
	result := dto.ToApplicationModel(updatedApplicationDTO, site.ID)

	// 5. Bind Site & Policies (Idempotent)
	err = a.bindSiteAndPolicies(result, site, policies)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *ApplicationServiceImpl) CompleteApplication(ctx context.Context, updatedApplication *model.Application) (*dto.ApplicationDTO, error) {
	// The application entity in SAC might contain additional attributes which are unknown or not related to this
	// operator. Instead of sending the updated application received from the operator, this function first fetch the
	// existing application in SAC and merge the updated application data to it in order not to override attributes
	// which have been updated in SAC but is not related here.
	foundApplicationDTO, err := a.sacClient.FindApplicationByID(*updatedApplication.ID)
	if err != nil {
		return nil, err
	}

	updatedApplicationDTO := dto.FromApplicationModel(updatedApplication)
	mergedApplicationDTO := dto.MergeApplication(foundApplicationDTO, updatedApplicationDTO)

	return mergedApplicationDTO, nil
}

func (a *ApplicationServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	log := logger.FromContext(ctx)
	log.Info("Deleting Application: '" + id.String() + "'...")

	err := a.sacClient.DeleteApplication(id)
	if err != nil {
		return err
	}

	log.Info("Application: '" + id.String() + "' deleted successfully.")
	return nil
}

func (a *ApplicationServiceImpl) bindSiteAndPolicies(
	application *model.Application,
	site *dto.SiteDTO,
	policies []dto.PolicyDTO,
) error {
	// Bind to Site (Idempotent)
	err := a.sacClient.BindApplicationToSite(*application.ID, site.ID)
	if err != nil {
		return typederror.WrapError(typederror.PartiallySuccessError, err)
	}

	// Attach Policies (Idempotent)
	var policyIds []uuid.UUID
	for _, policy := range policies {
		policyIds = append(policyIds, *policy.ID)
	}
	err = a.sacClient.UpdatePolicies(*application.ID, application.Type, policyIds)
	if err != nil {
		return typederror.WrapError(typederror.PartiallySuccessError, err)
	}

	return nil
}
