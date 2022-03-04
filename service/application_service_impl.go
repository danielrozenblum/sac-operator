package service

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/pkg/errors"

	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service/sac"
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"bitbucket.org/accezz-io/sac-operator/utils/typederror"
)

type ApplicationServiceImpl struct {
	sacClient sac.SecureAccessCloudClient
	log       logr.Logger
}

func (a *ApplicationServiceImpl) Reconcile(ctx context.Context, application *model.Application) (*ApplicationReconcileOutput, error) {

	output := &ApplicationReconcileOutput{}

	if application == nil {
		return output, fmt.Errorf("application cannot be nil, %w", typederror.UnrecoverableError)
	}

	if application.ToDelete {
		if application.ID == "" {
			return output, fmt.Errorf("application ID is nil, %w", typederror.UnrecoverableError)
		}
		err := a.delete(application.ID)
		if err != nil {
			return output, err
		}
		output.Deleted = true
		return output, nil
	}

	if application.ID == "" {
		applicationModel, err := a.create(application)
		output.SACApplicationID = applicationModel.ID
		return output, err
	}

	return nil, nil
}

func NewApplicationServiceImpl(sacClient sac.SecureAccessCloudClient, logger logr.Logger) ApplicationService {
	return &ApplicationServiceImpl{sacClient: sacClient, log: logger}
}

func (a *ApplicationServiceImpl) create(applicationToCreate *model.Application) (*model.Application, error) {
	a.log.Info("Trying to create application: " + applicationToCreate.String())

	// 1. Find Application by Name to verify the name isn't used
	appInSac, err := a.sacClient.FindApplicationByName(applicationToCreate.Name)
	if err != nil && err != sac.ErrorNotFound {
		return &model.Application{}, err
	}

	if appInSac.ID != "" {
		return &model.Application{}, fmt.Errorf("%w application %s already exist %s", typederror.UnrecoverableError, applicationToCreate.Name, appInSac.ID)
	}

	// 2. Validate Site Exists
	site, err := a.sacClient.FindSiteByName(applicationToCreate.Site)
	if err != nil {
		if errors.Is(err, sac.ErrorNotFound) {
			return &model.Application{}, fmt.Errorf("%w site %s does not exist", typederror.UnrecoverableError, applicationToCreate.Site)
		}
		return &model.Application{}, fmt.Errorf("%w could not validate site status", err)
	}

	// 3. Validate Policies Exists
	policies, err := a.sacClient.FindPoliciesByNames(append(applicationToCreate.AccessPolicies, applicationToCreate.ActivityPolicies...))
	if err != nil {
		if errors.Is(err, sac.ErrorNotFound) {
			return &model.Application{}, fmt.Errorf("%w policy does not exist %s", typederror.UnrecoverableError, err)
		}
		return &model.Application{}, fmt.Errorf("%w could not validate policy status", err)
	}

	// 4. createSite Application
	applicationDTO := dto.FromApplicationModel(applicationToCreate)
	createdApplicationDTO, err := a.sacClient.CreateApplication(applicationDTO)
	if err != nil {
		return &model.Application{}, err
	}

	// 5. Convert back to model
	result := dto.ToApplicationModel(createdApplicationDTO, site.ID)

	// 6. Bind Site & Policies (Idempotent)
	err = a.bindSiteToApplication(result, site)
	if err != nil {
		return &model.Application{
			ID: createdApplicationDTO.ID,
		}, err
	}

	err = a.bindPoliciesToApplication(result, policies)
	if err != nil {
		return &model.Application{
			ID: createdApplicationDTO.ID,
		}, err
	}

	a.log.Info("Application: '" + applicationToCreate.String() + "' created successfully.")
	return result, nil
}

func (a *ApplicationServiceImpl) update(updatedApplication *model.Application) (*model.Application, error) {
	a.log.Info("Dummy update Application")

	// 1. Validate Site Exists
	site, err := a.sacClient.FindSiteByName(updatedApplication.Site)
	if err != nil {
		return nil, fmt.Errorf("%w %s", typederror.UnrecoverableError, err)
	}

	// 2. Validate Policies Exists
	policies, err := a.sacClient.FindPoliciesByNames(append(updatedApplication.AccessPolicies, updatedApplication.ActivityPolicies...))
	if err != nil {
		return nil, fmt.Errorf("%w %s", typederror.UnrecoverableError, err)
	}

	// 3. update Application
	applicationDTO, err := a.completeApplication(updatedApplication)
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
	err = a.bindSiteToApplication(result, site)
	if err != nil {
		return nil, err
	}

	// 5. Bind Site & Policies (Idempotent)
	err = a.bindPoliciesToApplication(result, policies)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *ApplicationServiceImpl) completeApplication(updatedApplication *model.Application) (*dto.ApplicationDTO, error) {
	// The application entity in SAC might contain additional attributes which are unknown or not related to this
	// operator. Instead of sending the updated appli cation received from the operator, this function first fetch the
	// existing application in SAC and merge the updated application data to it in order not to override attributes
	// which have been updated in SAC but is not related here.
	foundApplicationDTO, err := a.sacClient.FindApplicationByID(updatedApplication.ID)
	if err != nil {
		return nil, err
	}

	updatedApplicationDTO := dto.FromApplicationModel(updatedApplication)
	mergedApplicationDTO := dto.MergeApplication(foundApplicationDTO, updatedApplicationDTO)

	return mergedApplicationDTO, nil
}

func (a *ApplicationServiceImpl) delete(id string) error {
	a.log.Info("Deleting Application: '" + id + "'...")

	err := a.sacClient.DeleteApplication(id)
	if err != nil {
		return err
	}

	a.log.Info("Application: '" + id + "' deleted successfully.")
	return nil
}

func (a *ApplicationServiceImpl) bindSiteToApplication(
	application *model.Application,
	site *dto.SiteDTO,
) error {
	// Bind to Site (Idempotent)
	err := a.sacClient.BindApplicationToSite(application.ID, site.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *ApplicationServiceImpl) bindPoliciesToApplication(
	application *model.Application,
	policies []dto.PolicyDTO,
) error {

	// Attach Policies (Idempotent)
	var policyIds []string
	for _, policy := range policies {
		policyIds = append(policyIds, policy.ID)
	}
	err := a.sacClient.UpdatePolicies(application.ID, application.Type, policyIds)
	if err != nil {
		return err
	}

	return nil
}
