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
		ID:                  uuid,
		TenantIdentifier:    site.Spec.TenantIdentifier,
		NumberOfConnectors:  site.Spec.NumberOfConnectors,
		ConnectorsNamespace: site.Spec.ConnectorsNamespace,
		EndpointURL:         site.Spec.EndpointURL,
	}
	connectors := []model.Connector{}

	for _, connector := range site.Status.Connectors {
		connectorUUID, err2 := utils.FromUIDType(connector.ConnectorID)
		if err2 != nil {
			return nil, err2
		}
		connectors = append(connectors, model.Connector{
			ConnectorID:     connectorUUID,
			Name:            connector.Name,
			ConnectorStatus: connector.ConnectorStatus,
		})
	}

	siteModel.Connectors = connectors

	return siteModel, nil
}

func (s *SiteConverter) UpdateStatus(siteModel *model.Site, site *accessv1.SiteStatus) error {

	var connectors []accessv1.Connector

	for i := range siteModel.Connectors {
		connectorID := utils.FromUUID(*siteModel.Connectors[i].ConnectorID)
		podID := utils.FromUUID(*siteModel.Connectors[i].ConnectorDeploymentID)
		connectors = append(connectors, accessv1.Connector{
			ConnectorID:     connectorID,
			PodID:           podID,
			Name:            siteModel.Connectors[i].Name,
			ConnectorStatus: siteModel.Connectors[i].ConnectorStatus,
		})
	}

	site.ID = utils.FromUUID(*siteModel.ID)
	site.Connectors = connectors

	return nil
}
