package service

import (
	"bitbucket.org/accezz-io/sac-operator/model"
	"context"
	"github.com/google/uuid"
	logger "sigs.k8s.io/controller-runtime/pkg/log"
)

type DummyApplicationServiceImpl struct{}

func NewDummyApplicationServiceImpl() ApplicationService {
	return &DummyApplicationServiceImpl{}
}

func (a *DummyApplicationServiceImpl) Create(ctx context.Context, applicationToCreate *model.Application) (*model.Application, error) {
	log := logger.FromContext(ctx)
	log.Info("Dummy Create Application")

	// 1. Create Application
	// 2. Bind to Site
	// 3. Attach Policies

	return model.NewApplicationBuilder().Build(), nil
}

func (a *DummyApplicationServiceImpl) Update(ctx context.Context, id uuid.UUID, updatedApplication *model.Application) (*model.Application, error) {
	log := logger.FromContext(ctx)
	log.Info("Dummy Update Application")
	return model.NewApplicationBuilder().WithID(id).Build(), nil
}

func (a *DummyApplicationServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	log := logger.FromContext(ctx)
	log.Info("Dummy Delete Application")
	return nil
}

func (a *DummyApplicationServiceImpl) GetById(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	log := logger.FromContext(ctx)
	log.Info("Dummy Get Application By Id")
	return model.NewApplicationBuilder().WithID(id).Build(), nil
}
