package service

//
//import (
//	"context"
//
//	"bitbucket.org/accezz-io/sac-operator/model"
//	"github.com/google/uuid"
//	logger "sigs.k8s.io/controller-runtime/pkg/log"
//)
//
//type DummyApplicationServiceImpl struct{}
//
//func NewDummyApplicationServiceImpl() ApplicationService {
//	return &DummyApplicationServiceImpl{}
//}
//
//func (a *DummyApplicationServiceImpl) Create(ctx context.Context, applicationToCreate *model.Application) (*model.Application, error) {
//	log := logger.FromContext(ctx)
//	log.Info("Dummy createSite Application")
//	return model.NewApplicationBuilder().Build(), nil
//}
//
//func (a *DummyApplicationServiceImpl) Update(ctx context.Context, updatedApplication *model.Application) (*model.Application, error) {
//	log := logger.FromContext(ctx)
//	log.Info("Dummy updateSiteAndPolicies Application")
//	return model.NewApplicationBuilder().WithID(updatedApplication.ID).Build(), nil
//}
//
//func (a *DummyApplicationServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
//	log := logger.FromContext(ctx)
//	log.Info("Dummy delete Application")
//	return nil
//}
//
//func (a *DummyApplicationServiceImpl) GetById(ctx context.Context, id uuid.UUID) (*model.Application, error) {
//	log := logger.FromContext(ctx)
//	log.Info("Dummy Get Application By Id")
//	return model.NewApplicationBuilder().WithID(&id).Build(), nil
//}
