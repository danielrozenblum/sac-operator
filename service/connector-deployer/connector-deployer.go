package connector_deployer

import (
	"context"

	"github.com/google/uuid"
)

type CreateConnectorInput struct {
	ConnectorID     *uuid.UUID
	SiteName        string
	Image           string
	Name            string
	EnvironmentVars map[string]string
}

type CreateConnctorOutput struct {
	DeploymentID *uuid.UUID
	Status       string
}

type GetConnectorsInput struct {
}

type ConnectorStatus string

const (
	OKConnectorStatus       = "OK"
	RecreateConnectorStatus = "Recreate"
	PendingConnectorStatus  = "Pending"
)

type Connector struct {
	ID     *uuid.UUID
	Status ConnectorStatus
}

type ConnnectorDeployer interface {
	CreateConnector(ctx context.Context, inputs *CreateConnectorInput) (*CreateConnctorOutput, error)
	GetConnectorsForSite(ctx context.Context, siteName string) ([]Connector, error)
}
