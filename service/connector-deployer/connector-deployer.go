package connector_deployer

import (
	"context"

	"github.com/google/uuid"
)

type DeployConnectorInput struct {
	ConnectorID     *uuid.UUID
	SiteName        string
	Image           string
	Name            string
	EnvironmentVars map[string]string
	Namespace       string
	SiteNamespace   string
}

type DeployConnectorOutput struct {
	DeploymentID *uuid.UUID
	Status       string
}

type ConnnectorDeployer interface {
	Deploy(ctx context.Context, inputs *DeployConnectorInput) (*DeployConnectorOutput, error)
}
