package model

type ConnectorConfiguration struct {
	ImagePullSecrets string
}

type Site struct {
	Name                   string
	SACSiteID              string
	NumberOfConnectors     int
	SiteNamespace          string
	ToDelete               bool
	ConnectorConfiguration *ConnectorConfiguration
}
