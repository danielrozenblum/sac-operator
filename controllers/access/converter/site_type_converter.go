package converter

import (
	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/utils"
)

// SiteConverter convert controller objects to service model
type SiteConverter struct{}

func NewSiteConverter() *SiteConverter {
	return &SiteConverter{}
}

func (s *SiteConverter) ConvertToServiceModel(site *accessv1.Site) (*model.Site, error) {
	uuid, err := utils.FromUIDType(site.Status.ID)
	if err != nil {
		return nil, err
	}

	siteModel := &model.Site{
		Name:                site.Name,
		SiteNamespace:       site.Namespace,
		SACSiteID:           uuid,
		TenantIdentifier:    site.Spec.TenantIdentifier,
		NumberOfConnectors:  site.Spec.NumberOfConnectors,
		ConnectorsNamespace: site.Spec.ConnectorsNamespace,
		EndpointURL:         site.Spec.EndpointURL,
	}
	connectors := []model.Connector{}

	for _, connector := range site.Status.Connectors {
		connectorUUID, err := utils.FromUIDType(connector.ConnectorID)
		if err != nil {
			return nil, err
		}
		deploymentUUID, err := utils.FromUIDType(connector.PodID)
		if err != nil {
			return nil, err
		}
		connectors = append(connectors, model.Connector{
			ConnectorID:           connectorUUID,
			ConnectorDeploymentID: deploymentUUID,
		})
	}

	siteModel.Connectors = connectors

	return siteModel, nil
}

func (s *SiteConverter) ConvertFromServiceModel(site *model.Site) *accessv1.SiteStatus {

	siteStatus := &accessv1.SiteStatus{}

	siteStatus.ID = utils.FromUUID(*site.SACSiteID)
	siteStatus.Connectors = make([]accessv1.Connector, len(site.Connectors))
	for i := range site.Connectors {
		siteStatus.Connectors[i].ConnectorID = utils.FromUUID(*site.Connectors[i].ConnectorID)
		siteStatus.Connectors[i].PodID = utils.FromUUID(*site.Connectors[i].ConnectorDeploymentID)
	}

	return siteStatus
}

func (s *SiteConverter) UpdateStatus(siteModel *model.Site, site *accessv1.SiteStatus) error {

	var connectors []accessv1.Connector

	for i := range siteModel.Connectors {
		connectorID := utils.FromUUID(*siteModel.Connectors[i].ConnectorID)
		podID := utils.FromUUID(*siteModel.Connectors[i].ConnectorDeploymentID)
		connectors = append(connectors, accessv1.Connector{
			ConnectorID: connectorID,
			PodID:       podID,
		})
	}

	site.ID = utils.FromUUID(*siteModel.SACSiteID)
	site.Connectors = connectors

	return nil
}
