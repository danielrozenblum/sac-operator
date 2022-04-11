package service

import (
	"context"

	"bitbucket.org/accezz-io/sac-operator/model"
)

//go:generate mockery -name ApplicationService -inpkg -case=underscore -output MockApplicationService
type ApplicationService interface {
	Reconcile(ctx context.Context, applicationToCreate *model.Application) (*ApplicationReconcileOutput, error)
}

type ApplicationReconcileOutput struct {
	Deleted          bool
	SACApplicationID string
}
