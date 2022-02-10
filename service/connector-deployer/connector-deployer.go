package connector_deployer

import (
	"context"

	accessv1 "bitbucket.org/accezz-io/sac-operator/apis/access/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/uuid"
)

type Owner struct {
	object *v1.Object
}

type DeployConnectorInput struct {
	ConnectorID     *uuid.UUID
	SiteName        string
	Image           string
	Name            string
	EnvironmentVars map[string]string
	Site            *accessv1.Site
	Namespace       string
}

type DeployConnectorOutput struct {
	DeploymentID *uuid.UUID
	Status       string
}

type ConnnectorDeployer interface {
	Deploy(ctx context.Context, inputs *DeployConnectorInput) (*DeployConnectorOutput, error)
}
