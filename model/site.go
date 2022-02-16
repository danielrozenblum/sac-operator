package model

type Site struct {
	Name                string
	SACSiteID           string
	TenantIdentifier    string
	NumberOfConnectors  int
	EndpointURL         string
	ConnectorsNamespace string
	SiteNamespace       string
	ToDelete            bool
	Deleted             bool
}
