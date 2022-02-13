package service

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/go-logr/logr"

	connector_deployer "bitbucket.org/accezz-io/sac-operator/service/connector-deployer"

	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service/sac"
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"github.com/google/uuid"
)

type SiteServiceImpl struct {
	connectorDeployer connector_deployer.ConnnectorDeployer
	sacClient         sac.SecureAccessCloudClient
	log               logr.Logger
}

func NewSiteServiceImpl(sacClient sac.SecureAccessCloudClient,
	connectorDeployer connector_deployer.ConnnectorDeployer,
	log logr.Logger,
) *SiteServiceImpl {
	return &SiteServiceImpl{
		sacClient:         sacClient,
		connectorDeployer: connectorDeployer,
		log:               log,
	}
}

func (s *SiteServiceImpl) Reconcile(ctx context.Context, site *model.Site) error {

	var err error
	// in case site does not exist in SAC
	if site.SACSiteID == nil {
		err = s.createSite(ctx, site)
		if err != nil {
			return err
		}
	}

	connectors, err := s.connectorDeployer.GetConnectorsForSite(ctx, site.Name)
	if err != nil {
		return err
	}
	//not enough deployed connectors
	for i := len(connectors); i < site.NumberOfConnectors; i++ {
		err = s.createConnector(ctx, site)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *SiteServiceImpl) getDeployConnectorInputs(connector *dto.Connector, site *model.Site) (*connector_deployer.CreateConnectorInput, error) {
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

func (s *SiteServiceImpl) createSite(ctx context.Context, site *model.Site) error {
	sacSite := dto.FromSiteModel(site)
	siteDto, err := s.sacClient.CreateSite(sacSite)
	if err != nil {
		return err
	}

	site.SACSiteID = siteDto.ID

	for i := 0; i < site.NumberOfConnectors; i++ {

		connector, err := s.sacClient.CreateConnector(siteDto, s.getConnectorName(site))
		if err != nil {
			return err
		}
		deployConnectorInput, err := s.getDeployConnectorInputs(connector, site)
		if err != nil {
			return err
		}

		deployConnectorOutput, err := s.connectorDeployer.CreateConnector(ctx, deployConnectorInput)
		if err != nil {
			return err
		}

		site.Connectors = append(site.Connectors, model.Connector{
			ConnectorID:           deployConnectorInput.ConnectorID,
			ConnectorDeploymentID: deployConnectorOutput.DeploymentID,
			Name:                  deployConnectorInput.Name,
		})
	}

	return nil

}

func (s *SiteServiceImpl) createConnector(ctx context.Context, site *model.Site) error {

	s.log.WithValues("site name", site.Name).Info("creating connector in site")

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

	deployConnectorOutput, err := s.connectorDeployer.CreateConnector(ctx, deployConnectorInput)
	if err != nil {
		return err
	}
	s.log.WithValues("id", deployConnectorOutput.DeploymentID).Info("deployed new connector")

	site.Connectors = append(site.Connectors, model.Connector{
		ConnectorID:           deployConnectorInput.ConnectorID,
		ConnectorDeploymentID: deployConnectorOutput.DeploymentID,
		Name:                  deployConnectorInput.Name,
	})

	return nil

}

func (s *SiteServiceImpl) getConnectorName(site *model.Site) string {

	return fmt.Sprintf("%s-%s-%s", site.Name, site.SiteNamespace, rand.String(4))

}
