package connector_deployer

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateConnectorInput struct {
	ConnectorID     *uuid.UUID
	SiteName        string
	Image           string
	Name            string
	EnvironmentVars map[string]string
}

type GetConnectorsInput struct {
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

type ConnnectorDeployer interface {
	CreateConnector(ctx context.Context, inputs *CreateConnectorInput) error
	DeleteConnector(ctx context.Context, Name string) error
	GetConnectorsForSite(ctx context.Context, siteName string) ([]Connector, error)
}
