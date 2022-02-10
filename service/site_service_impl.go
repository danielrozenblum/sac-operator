package service

import (
	"context"
	"errors"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"

	connector_deployer "bitbucket.org/accezz-io/sac-operator/service/connector-deployer"

	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service/sac"
	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"
	"github.com/google/uuid"
)

type SiteServiceImpl struct {
	connectorDeployer connector_deployer.ConnnectorDeployer
	sacClient         sac.SecureAccessCloudClient
}

func NewSiteServiceImpl(sacClient sac.SecureAccessCloudClient, connectorDeployer connector_deployer.ConnnectorDeployer) *SiteServiceImpl {
	return &SiteServiceImpl{
		sacClient:         sacClient,
		connectorDeployer: connectorDeployer,
	}
}

func (s *SiteServiceImpl) isSiteExist(ctx context.Context, name string) (bool, error) {
	_, err := s.sacClient.FindSiteByName(name)
	if err != nil {
		switch {
		case errors.Is(err, sac.ErrorNotFound):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (s *SiteServiceImpl) Create(ctx context.Context, site *model.Site, siteCRD *accessv1.Site) error {

	sacSite := dto.FromSiteModel(site)
	siteDto, err := s.sacClient.CreateSite(sacSite)
	if err != nil {
		return err
	}

	site.ID = siteDto.ID

	for i := 0; i < site.NumberOfConnectors; i++ {
		connector, err := s.sacClient.CreateConnector(siteDto)
		if err != nil {
			return err
		}
		deployConnectorInput, err := s.getDeployConnectorInputs(connector, site, siteCRD)
		if err != nil {
			return err
		}

		deployConnectorOutput, err := s.connectorDeployer.Deploy(ctx, deployConnectorInput)
		if err != nil {
			return err
		}

		site.Connectors = append(site.Connectors, model.Connector{
			ConnectorID:           deployConnectorInput.ConnectorID,
			ConnectorDeploymentID: deployConnectorOutput.DeploymentID,
			Name:                  deployConnectorInput.Name,
			ConnectorStatus:       deployConnectorOutput.Status,
		})
	}

	return nil

}

func (s *SiteServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return s.sacClient.DeleteSite(id)
}

func (s *SiteServiceImpl) getDeployConnectorInputs(connector *dto.Connector, site *model.Site, siteCRD *accessv1.Site) (*connector_deployer.DeployConnectorInput, error) {
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

	return &connector_deployer.DeployConnectorInput{
		ConnectorID:     &ConnectorID,
		SiteName:        site.Name,
		Image:           "luminate/connector:2.10.1", //TODO: waiting for https://jira.luminate.io/browse/AC-27957
		Name:            connector.Name,
		EnvironmentVars: envs,
		Namespace:       site.ConnectorsNamespace,
		Site:            siteCRD,
	}, nil
}
