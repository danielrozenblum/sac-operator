package model

import (
	"github.com/google/uuid"
)

type ConnectorDeploymentStatus string

const (
	RunningConnectorStatus ConnectorDeploymentStatus = "Running"
)

type Connector struct {
	ConnectorID           *uuid.UUID
	ConnectorDeploymentID *uuid.UUID
	Version               string
	Name                  string
}

type Site struct {
	Name                string
	ID                  *uuid.UUID
	TenantIdentifier    string
	NumberOfConnectors  int
	EndpointURL         string
	Connectors          []Connector
	ConnectorsNamespace string
	SiteNamespace       string
}
