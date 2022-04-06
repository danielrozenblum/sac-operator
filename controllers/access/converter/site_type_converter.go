package converter

import (
	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"
	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service"
)

// SiteConverter convert controller objects to service model
type SiteConverter struct{}

func NewSiteConverter() *SiteConverter {
	return &SiteConverter{}
}

func (s *SiteConverter) ConvertToServiceModel(site *accessv1.Site) *model.Site {

	connectorConfiguration := &model.ConnectorConfiguration{}
	if site.Spec.ImagePullSecret != "" {
		connectorConfiguration.ImagePullSecrets = site.Spec.ImagePullSecret
	}

	siteModel := &model.Site{
		Name:                   site.Name,
		SiteNamespace:          site.Namespace,
		SACSiteID:              site.Status.ID,
		NumberOfConnectors:     site.Spec.NumberOfConnectors,
		ToDelete:               !site.ObjectMeta.DeletionTimestamp.IsZero(),
		ConnectorConfiguration: connectorConfiguration,
	}

	return siteModel
}

func (s *SiteConverter) ConvertFromServiceOutput(site *service.SiteReconcileOutput) accessv1.SiteStatus {

	siteStatus := accessv1.SiteStatus{
		ID:                        "",
		HealthyConnectors:         map[string]string{},
		UnHealthyConnectors:       map[string]string{},
		NumberOfHealthyConnectors: 0,
	}

	siteStatus.ID = site.SACSiteID
	for i := range site.HealthyConnectors {
		siteStatus.HealthyConnectors[site.HealthyConnectors[i].DeploymentName] = site.HealthyConnectors[i].SacID
	}
	for i := range site.UnHealthyConnectors {
		siteStatus.UnHealthyConnectors[site.UnHealthyConnectors[i].DeploymentName] = site.UnHealthyConnectors[i].SacID
	}
	siteStatus.NumberOfHealthyConnectors = len(site.HealthyConnectors)

	return siteStatus
}

func (s *SiteConverter) UpdateStatus(siteModel *model.Site, site *accessv1.SiteStatus) error {

	site.ID = siteModel.SACSiteID

	return nil
}
