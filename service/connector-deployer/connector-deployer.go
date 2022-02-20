package connector_deployer

import (
	"context"
	"time"
)

type CreateConnectorInput struct {
	ConnectorID     string
	SiteName        string
	Image           string
	Name            string
	EnvironmentVars map[string]string
}

type ConnectorStatus string

const (
	OKConnectorStatus       = "OK"
	ToDeleteConnectorStatus = "ToDelete"
)

type Connector struct {
	DeploymentName   string
	SACID            string
	Status           ConnectorStatus
	CreatedTimeStamp time.Time
}

//go:generate mockery --name=ConnectorDeployer --inpackage --case=underscore --output=mockConnectorDeployerInterface
type ConnectorDeployer interface {
	CreateConnector(ctx context.Context, inputs *CreateConnectorInput) (string, error)
	DeleteConnector(ctx context.Context, Name string) error
	GetConnectorsForSite(ctx context.Context, siteName string) ([]Connector, error)
}
