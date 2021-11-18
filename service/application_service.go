package service

import (
	"bitbucket.org/accezz-io/sac-operator/model"
	"context"
	"github.com/google/uuid"
)

//go:generate mockery -name ApplicationService -inpkg -case=underscore -output MockApplicationService
type ApplicationService interface {
	Create(ctx context.Context, applicationToCreate *model.Application) (*model.Application, error)
	Update(ctx context.Context, id uuid.UUID, updatedApplication *model.Application) (*model.Application, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (*model.Application, error)
}
