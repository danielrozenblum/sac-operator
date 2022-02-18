package service

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"bitbucket.org/accezz-io/sac-operator/service/repository"

	"bitbucket.org/accezz-io/sac-operator/utils"

	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/go-logr/logr"

	connector_deployer "bitbucket.org/accezz-io/sac-operator/service/connector-deployer"

	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service/sac"
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"github.com/google/uuid"
)

var UnrecoverableError = fmt.Errorf("unrecoverable error")

type SiteServiceImpl struct {
	connectorDeployer connector_deployer.ConnectorDeployer
	sacClient         sac.SecureAccessCloudClient
	repo              repository.Repository
	log               logr.Logger
}

func NewSiteServiceImpl(sacClient sac.SecureAccessCloudClient,
	connectorDeployer connector_deployer.ConnectorDeployer,
	log logr.Logger,
	repo repository.Repository,
) *SiteServiceImpl {
	return &SiteServiceImpl{
		sacClient:         sacClient,
		connectorDeployer: connectorDeployer,
		log:               log,
		repo:              repo,
	}
}

func (s *SiteServiceImpl) Reconcile(ctx context.Context, site *model.Site) (*SiteReconcileOutput, error) {
	output := &SiteReconcileOutput{}

	if site == nil {
		s.log.Error(fmt.Errorf(""), "site model cannot be nil")
		return output, UnrecoverableError
	}

	if site.ToDelete {
		err := s.deleteSiteInSAC(ctx, site, output)
		return output, err // nothing to reconcile other than deleting the site in SAC
	}

	if site.SACSiteID == "" {
		err := s.createSiteInSAC(ctx, site, output)
		if err != nil {
			return output, err
		}
	} else {
		output.SACSiteID = site.SACSiteID
	}

	connectors, err := s.connectorDeployer.GetConnectorsForSite(ctx, site.Name)
	if err != nil {
		return output, err
	}

	var healthyConnectors []connector_deployer.Connector
	var toDeleteConnectors []connector_deployer.Connector

	for i := range connectors {
		switch connectors[i].Status {
		case connector_deployer.ToDeleteConnectorStatus:
			toDeleteConnectors = append(toDeleteConnectors, connectors[i])
			output.UnHealthyConnectors[connectors[i].DeploymentName] = connectors[i].SACID
		case connector_deployer.OKConnectorStatus:
			healthyConnectors = append(healthyConnectors, connectors[i])
			output.HealthConnectors[connectors[i].DeploymentName] = connectors[i].SACID
		}
	}

	s.log.WithValues("site", site.Name,
		"desired", site.NumberOfConnectors,
		"healthyConnectors", len(healthyConnectors),
		"toDeleteConnectors", len(toDeleteConnectors)).
		Info("connectors status for site")

	err = s.reconcileToDeleteConnectors(ctx, site, toDeleteConnectors)
	if err != nil {
		return output, err
	}

	//var err error
	//err = s.reconcileSiteInSAC(ctx, site)
	//if err != nil {
	//	return s.handleOutput(site, err)
	//}
	//if site.ToDelete {
	//	return &SiteReconcileOutput{}, nil
	//}
	//
	//connectors, err := s.connectorDeployer.GetConnectorsForSite(ctx, site.Name)
	//if err != nil {
	//	return s.handleOutput(site, err)
	//}
	//
	//var healthyConnectors []connector_deployer.Connector
	//var toDeleteConnectors []connector_deployer.Connector
	//
	//for i := range connectors {
	//	switch connectors[i].Status {
	//	case connector_deployer.ToDeleteConnectorStatus:
	//		toDeleteConnectors = append(toDeleteConnectors, connectors[i])
	//	case connector_deployer.OKConnectorStatus:
	//		healthyConnectors = append(healthyConnectors, connectors[i])
	//	}
	//}
	//
	//s.log.WithValues("site", site.Name, "desired", site.NumberOfConnectors, "healthyConnectors", len(healthyConnectors), "toDeleteConnectors", len(toDeleteConnectors)).Info("connectors")
	//
	//err = s.reconcileToDeleteConnectors(ctx, site, toDeleteConnectors)
	//if err != nil {
	//	return s.handleOutput(site, err)
	//}
	//
	//err = s.reconcileDesiredNumberOfConnectors(ctx, site, healthyConnectors)
	//if err != nil {
	//	return s.handleOutput(site, err)
	//}
	//
	//return s.handleOutput(site, err)

	return output, nil
}

func (s *SiteServiceImpl) createSiteInSAC(ctx context.Context, site *model.Site, output *SiteReconcileOutput) error {

	sacSite := dto.FromSiteModel(site)
	siteDto, err := s.sacClient.CreateSite(sacSite)
	if err != nil {
		if sac.IsConflict(err) {
			return fmt.Errorf("%w site already exist", UnrecoverableError)
		}
		return err
	}

	output.SACSiteID = siteDto.ID

	return nil

}

func (s *SiteServiceImpl) deleteSiteInSAC(ctx context.Context, site *model.Site, output *SiteReconcileOutput) error {

	err := s.sacClient.DeleteSite(site.SACSiteID)
	if err != nil {
		return err
	}

	output.Deleted = true

	return nil

}

func (s *SiteServiceImpl) reconcileToDeleteConnectors(ctx context.Context, site *model.Site, connectors []connector_deployer.Connector) error {

	for i := range connectors {
		s.log.WithValues(
			"sac id", connectors[i].SACID,
			"pod name", connectors[i].DeploymentName,
		).Info("deleting unhealthy connector",
			connectors[i].SACID,
			connectors[i].DeploymentName,
		)
		err := s.deleteConnector(ctx, connectors[i].SACID, connectors[i].DeploymentName)
		if err != nil {
			return err
		}
		connectors = append(connectors[:i], connectors[i+1:]...) // if delete ended successfully, delete it from connectors list
	}

	return nil

}

func (s *SiteServiceImpl) reconcileDesiredNumberOfConnectors(ctx context.Context, site *model.Site, connectors []connector_deployer.Connector) error {

	if site.ToDelete {
		return nil
	}

	s.log.WithValues("site name", site.Name, "number of connectors", len(connectors), "desired number of connectors", site.NumberOfConnectors).Info("actual vs desired number of connectors")

	if site.NumberOfConnectors == len(connectors) {
		return nil
	}

	if site.NumberOfConnectors > len(connectors) {
		for i := len(connectors); i < site.NumberOfConnectors; i++ {
			err := s.createConnector(ctx, site)
			if err != nil {
				return err
			}
		}
		return nil
	}

	sort.Slice(connectors, func(i, j int) bool {
		return connectors[i].CreatedTimeStamp.Before(connectors[j].CreatedTimeStamp)
	})

	if site.NumberOfConnectors < len(connectors) {
		s.log.Info("deleting oldest connectors")
		for i := len(connectors) - 1; i >= site.NumberOfConnectors; i-- {
			err := s.deleteConnector(ctx, connectors[i].SACID, connectors[i].DeploymentName)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func (s *SiteServiceImpl) reconcilerDanglingConnectorsFromSAC(ctx context.Context, site *model.Site) error {

	if site.ToDelete {
		return nil
	}
	connectors, err := s.connectorDeployer.GetConnectorsForSite(ctx, site.Name)

	if err != nil {
		return err
	}

	podIDs := func() []string {
		var ids []string
		for i := range connectors {
			ids = append(ids, connectors[i].SACID)
		}
		return ids
	}()

	// removing dangling from sac
	sacListOfConnectors, err := s.sacClient.ListConnectorsBySite(site.Name)
	if err != nil {
		return err
	}

	toDelete := utils.Subtruct(sacListOfConnectors, podIDs)

	for i := range toDelete {
		s.log.Info("deleting sac connector", "uuid", toDelete[i])
		err = s.sacClient.DeleteConnector(toDelete[i])
		if err != nil {
			s.log.Error(err, "could not delete sac connector", "uuid", toDelete[i])
		}
	}

	return nil

}

func (s *SiteServiceImpl) createConnector(ctx context.Context, site *model.Site) error {

	siteDto, err := s.sacClient.FindSiteByName(site.Name)
	if err != nil {
		return err
	}

	connector, err := s.sacClient.CreateConnector(siteDto, s.getConnectorName(site))
	if err != nil {
		return err
	}

	s.log.WithValues("id", connector.ID, "name", connector.Name).Info("created connector in sac")
	deployConnectorInput, err := s.getDeployConnectorInputs(connector, site)
	if err != nil {
		return err
	}

	err = s.connectorDeployer.CreateConnector(ctx, deployConnectorInput)
	if err != nil {
		return err
	}
	s.log.Info("deployed new connector")

	return nil

}

func (s *SiteServiceImpl) getDeployConnectorInputs(connector *dto.ConnectorObjects, site *model.Site) (*connector_deployer.CreateConnectorInput, error) {
	ConnectorID, err := uuid.Parse(connector.ID)
	if err != nil {
		return nil, err
	}

	envs := map[string]string{
		"ENDPOINT_URL":           site.EndpointURL,
		"TENANT_IDENTIFIER":      site.TenantIdentifier,
		"HTTPS_SKIP_CERT_VERIFY": "true",
		"OTP":                    connector.Otp,
	}

	return &connector_deployer.CreateConnectorInput{
		ConnectorID:     &ConnectorID,
		SiteName:        site.Name,
		Image:           "luminate/connector:2.10.1", //TODO: waiting for https://jira.luminate.io/browse/AC-27957
		Name:            connector.Name,
		EnvironmentVars: envs,
	}, nil
}

func (s *SiteServiceImpl) getConnectorName(site *model.Site) string {

	return fmt.Sprintf("%s-%s-%s", site.Name, site.SiteNamespace, rand.String(4))

}

func (s *SiteServiceImpl) deleteConnector(ctx context.Context, sacID, podID string) error {

	err := s.sacClient.DeleteConnector(sacID)
	if err != nil {
		return err
	}

	err = s.connectorDeployer.DeleteConnector(ctx, podID)
	if err != nil {
		return err
	}

	return nil

}

func (s *SiteServiceImpl) handleOutput(site *model.Site, err error) (*SiteReconcileOutput, error) {

	if err == nil {
		return &SiteReconcileOutput{}, nil
	}

	if errors.Is(err, UnrecoverableError) {
		s.log.WithValues("site name", site.Name).Error(err, "unrecoverable error")
		return &SiteReconcileOutput{}, nil
	}

	return &SiteReconcileOutput{}, err

}
