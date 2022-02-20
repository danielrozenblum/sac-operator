package service

import (
	"context"
	"fmt"

	"bitbucket.org/accezz-io/sac-operator/utils"

	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/go-logr/logr"

	connector_deployer "bitbucket.org/accezz-io/sac-operator/service/connector-deployer"

	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service/sac"
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
)

var UnrecoverableError = fmt.Errorf("unrecoverable error")

type SiteServiceImpl struct {
	connectorDeployer connector_deployer.ConnectorDeployer
	sacClient         sac.SecureAccessCloudClient
	log               logr.Logger
}

func NewSiteServiceImpl(sacClient sac.SecureAccessCloudClient,
	connectorDeployer connector_deployer.ConnectorDeployer,
	log logr.Logger,
) *SiteServiceImpl {
	return &SiteServiceImpl{
		sacClient:         sacClient,
		connectorDeployer: connectorDeployer,
		log:               log,
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

	for i := range connectors {
		switch connectors[i].Status {
		case connector_deployer.ToDeleteConnectorStatus:
			output.UnHealthyConnectors = append(output.UnHealthyConnectors, Connectors{
				CreatedTimestamp: connectors[i].CreatedTimeStamp,
				DeploymentName:   connectors[i].DeploymentName,
				SacID:            connectors[i].SACID,
			})
		case connector_deployer.OKConnectorStatus:
			output.HealthyConnectors = append(output.HealthyConnectors, Connectors{
				CreatedTimestamp: connectors[i].CreatedTimeStamp,
				DeploymentName:   connectors[i].DeploymentName,
				SacID:            connectors[i].SACID,
			})
		}
	}

	s.log.WithValues("site", site.Name,
		"desired", site.NumberOfConnectors,
		"healthyConnectors", len(output.HealthyConnectors),
		"toDeleteConnectors", len(output.UnHealthyConnectors)).
		Info("connectors status for site")

	if len(output.UnHealthyConnectors) > 0 {
		for i := range output.UnHealthyConnectors {
			err = s.deleteConnector(ctx, output.UnHealthyConnectors[i].SacID, output.UnHealthyConnectors[i].DeploymentName)
			if err != nil {
				return output, err
			}
			output.UnHealthyConnectors = append(output.UnHealthyConnectors[:i], output.UnHealthyConnectors[i+1:]...)
		}
	}

	switch {
	case len(output.HealthyConnectors) == site.NumberOfConnectors:
		break
	case len(output.HealthyConnectors) < site.NumberOfConnectors:
		for i := len(output.HealthyConnectors); i < site.NumberOfConnectors; i++ {
			connector, err := s.createConnector(ctx, site)
			if err != nil {
				return output, err
			}
			output.HealthyConnectors = append(output.HealthyConnectors, connector)
		}
	case len(output.HealthyConnectors) > site.NumberOfConnectors:
		sortConnectorsBtCreatedTimestamp(output.HealthyConnectors)
		for i := len(output.HealthyConnectors); i > site.NumberOfConnectors; i-- {
			err = s.deleteConnector(ctx, output.HealthyConnectors[i].SacID, output.HealthyConnectors[i].DeploymentName)
			if err != nil {
				return output, err
			}
			output.HealthyConnectors = append(output.HealthyConnectors[:i], output.HealthyConnectors[i+1:]...)
		}
	}

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

func (s *SiteServiceImpl) createConnector(ctx context.Context, site *model.Site) (Connectors, error) {

	connector := Connectors{}

	siteDto, err := s.sacClient.FindSiteByName(site.Name)
	if err != nil {
		return connector, err
	}

	sacConnector, err := s.sacClient.CreateConnector(siteDto, s.getConnectorName(site))
	if err != nil {
		return connector, err
	}
	connector.SacID = sacConnector.ID
	s.log.WithValues("id", sacConnector.ID, "name", sacConnector.Name).Info("created connector in sac")

	deployConnectorInput := s.getDeployConnectorInputs(sacConnector, site)

	deploymentName, err := s.connectorDeployer.CreateConnector(ctx, deployConnectorInput)
	if err != nil {
		return connector, err
	}
	connector.DeploymentName = deploymentName
	s.log.Info("deployed new connector")

	return connector, nil

}

func (s *SiteServiceImpl) getDeployConnectorInputs(connector *dto.ConnectorObjects, site *model.Site) *connector_deployer.CreateConnectorInput {

	envs := map[string]string{
		"ENDPOINT_URL":           site.EndpointURL,
		"TENANT_IDENTIFIER":      site.TenantIdentifier,
		"HTTPS_SKIP_CERT_VERIFY": "true",
		"OTP":                    connector.Otp,
	}

	return &connector_deployer.CreateConnectorInput{
		ConnectorID:     connector.ID,
		SiteName:        site.Name,
		Image:           "luminate/connector:2.10.1", //TODO: waiting for https://jira.luminate.io/browse/AC-27957
		Name:            connector.Name,
		EnvironmentVars: envs,
	}
}

func (s *SiteServiceImpl) getConnectorName(site *model.Site) string {

	return fmt.Sprintf("%s-%s-%s", site.Name, site.SiteNamespace, rand.String(4))

}

func (s *SiteServiceImpl) deleteConnector(ctx context.Context, sacID, podName string) error {

	err := s.sacClient.DeleteConnector(sacID)
	if err != nil {
		return err
	}

	err = s.connectorDeployer.DeleteConnector(ctx, podName)
	if err != nil {
		return err
	}

	return nil

}
