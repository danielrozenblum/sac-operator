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

	err := a.getSiteAndPoliciesIDs(application)
	if err != nil {
		return output, err
	}

	if application.ID == "" {
		err = a.create(application)
		if err != nil {
			return output, err
		}
	} else {
		err = a.updateApplication(application)
		if err != nil {
			output.SACApplicationID = application.ID
			return output, err
		}
	}

	output.SACApplicationID = application.ID

	err = a.updateSiteAndPolicies(application)
	if err != nil {
		return output, err
	}

	output.SiteID = application.SiteId
	output.PoliciesIDs = application.PoliciesIDS

	return output, nil
}

func NewApplicationServiceImpl(sacClient sac.SecureAccessCloudClient, logger logr.Logger) ApplicationService {
	return &ApplicationServiceImpl{sacClient: sacClient, log: logger}
}

func (a *ApplicationServiceImpl) create(applicationToCreate *model.Application) error {
	a.log.Info("creating application: " + applicationToCreate.String())

	// 1. Find Application by Name to verify the name isn't used
	appInSac, err := a.sacClient.FindApplicationByName(applicationToCreate.Name)
	if err != nil && err != sac.ErrorNotFound {
		return err
	}

	if appInSac.ID != "" {
		return fmt.Errorf("%w application %s already exist %s", typederror.UnrecoverableError, applicationToCreate.Name, appInSac.ID)
	}

	// 4. create Application
	applicationDTO := dto.FromApplicationModel(applicationToCreate)
	createdApplicationDTO, err := a.sacClient.CreateApplication(applicationDTO)
	if err != nil {
		return err
	}
	applicationToCreate.ID = createdApplicationDTO.ID

	return nil
}

func (a *ApplicationServiceImpl) getSiteAndPoliciesIDs(applicationToCreate *model.Application) error {

	// 2. Validate SiteName Exists
	site, err := a.sacClient.FindSiteByName(applicationToCreate.SiteName)
	if err != nil {
		if errors.Is(err, sac.ErrorNotFound) {
			return fmt.Errorf("%w site %s does not exist", typederror.UnrecoverableError, applicationToCreate.SiteName)
		}
		return fmt.Errorf("%w could not validate site status", err)
	}

	applicationToCreate.SiteId = site.ID

	// 3. Validate Policies Exists
	policies, err := a.sacClient.FindPoliciesByNames(append(applicationToCreate.AccessPoliciesNames, applicationToCreate.ActivityPoliciesNames...))
	if err != nil {
		if errors.Is(err, sac.ErrorNotFound) {
			return fmt.Errorf("%w policy does not exist %s", typederror.UnrecoverableError, err)
		}
		return fmt.Errorf("%w could not validate policy status", err)
	}

	for i := range policies {
		applicationToCreate.PoliciesIDS = append(applicationToCreate.PoliciesIDS, policies[i].ID)
	}

	return nil
}

func (a *ApplicationServiceImpl) updateSiteAndPolicies(application *model.Application) error {

	// 5. Bind SiteName & Policies (Idempotent)
	err := a.bindSiteToApplication(application)
	if err != nil {
		return err
	}

	// 5. Bind SiteName & Policies (Idempotent)
	err = a.bindPoliciesToApplication(application)
	if err != nil {
		return err
	}

	return nil
}

func (a *ApplicationServiceImpl) updateApplication(application *model.Application) error {

	// 3. updateSiteAndPolicies Application
	applicationDTO := dto.FromApplicationModel(application)

	_, err := a.sacClient.UpdateApplication(applicationDTO)
	if err != nil {
		if errors.Is(err, sac.ErrorNotFound) {
			return fmt.Errorf("%w application id %s not found", typederror.UnrecoverableError, application.ID)
		}
		return err
	}

	return nil
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
) error {
	// Bind to SiteName (Idempotent)
	err := a.sacClient.BindApplicationToSite(application.ID, application.SiteId)
	if err != nil {
		return err
	}

	return nil
}

func (a *ApplicationServiceImpl) bindPoliciesToApplication(
	application *model.Application,
) error {

	err := a.sacClient.UpdatePolicies(application.ID, application.Type, application.PoliciesIDS)
	if err != nil {
		return err
	}

	return nil
}
