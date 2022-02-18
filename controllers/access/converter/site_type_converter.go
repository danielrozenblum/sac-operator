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

	siteModel := &model.Site{
		Name:                site.Name,
		SiteNamespace:       site.Namespace,
		SACSiteID:           site.Status.ID,
		TenantIdentifier:    site.Spec.TenantIdentifier,
		NumberOfConnectors:  site.Spec.NumberOfConnectors,
		ConnectorsNamespace: site.Spec.ConnectorsNamespace,
		EndpointURL:         site.Spec.EndpointURL,
		ToDelete:            !site.ObjectMeta.DeletionTimestamp.IsZero(),
	}

	return siteModel
}

func (s *SiteConverter) ConvertFromServiceOutput(site *service.SiteReconcileOutput) accessv1.SiteStatus {

	siteStatus := accessv1.SiteStatus{}

	siteStatus.ID = site.SACSiteID

	return siteStatus
}

func (s *SiteConverter) UpdateStatus(siteModel *model.Site, site *accessv1.SiteStatus) error {

	site.ID = siteModel.SACSiteID

	return nil
}
