package service

import (
	"context"
	"testing"

	"bitbucket.org/accezz-io/sac-operator/utils/typederror"

	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"

	"github.com/stretchr/testify/mock"

	"bitbucket.org/accezz-io/sac-operator/model"
	"bitbucket.org/accezz-io/sac-operator/service/sac"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestApplicationServiceImpl_Reconcile(t *testing.T) {
	errorFromSacService := typederror.UnknownError
	tests := []struct {
		name      string
		setupFunc func() (ApplicationService, *model.Application)
		output    *ApplicationReconcileOutput
		err       error
	}{
		{
			name: "nil application flow",
			setupFunc: func() (ApplicationService, *model.Application) {
				sacClient := &sac.MockSecureAccessCloudClient{}
				testLog := ctrl.Log.WithName("test")
				return NewApplicationServiceImpl(sacClient, testLog), nil
			},
			output: &ApplicationReconcileOutput{},
			err:    typederror.UnrecoverableError,
		},
		{
			name: "[delete application flow] no application ID",
			setupFunc: func() (ApplicationService, *model.Application) {
				sacClient := &sac.MockSecureAccessCloudClient{}
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					ToDelete: true,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    typederror.UnrecoverableError,
		},
		{
			name: "[delete application flow] failed to delete",
			setupFunc: func() (ApplicationService, *model.Application) {
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("DeleteApplication", mock.AnythingOfType("string")).Return(errorFromSacService)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					ToDelete: true,
					ID:       "uuid",
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{
				Deleted: false,
			},
			err: errorFromSacService,
		},
		{
			name: "[delete application flow] success flow",
			setupFunc: func() (ApplicationService, *model.Application) {
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("DeleteApplication", mock.AnythingOfType("string")).Return(nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					ToDelete: true,
					ID:       "uuid",
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{
				Deleted: true,
			},
			err: nil,
		},
		{
			name: "[new application flow] application already exist",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{ID: "uuid"}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name: appName,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    typederror.UnrecoverableError,
		},
		{
			name: "[new application flow] find application error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, errorFromSacService)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name: appName,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    errorFromSacService,
		},
		{
			name: "[new application flow] find application site does not exist ",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{}, sac.ErrorNotFound)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name: appName,
					Site: siteToUse,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    typederror.UnrecoverableError,
		},
		{
			name: "[new application flow] find site by name error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{}, errorFromSacService)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name: appName,
					Site: siteToUse,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    errorFromSacService,
		},
		{
			name: "[new application flow] find policy by name not found error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				accessPolicy := []string{"1", "2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{}, nil)
				sacClient.On("FindPoliciesByNames", accessPolicy).Return([]dto.PolicyDTO{}, sac.ErrorNotFound)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:           appName,
					Site:           siteToUse,
					AccessPolicies: accessPolicy,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    typederror.UnrecoverableError,
		},
		{
			name: "[new application flow] find policy by name error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				accessPolicy := []string{"1", "2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{}, nil)
				sacClient.On("FindPoliciesByNames", accessPolicy).Return([]dto.PolicyDTO{}, errorFromSacService)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:           appName,
					Site:           siteToUse,
					AccessPolicies: accessPolicy,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    errorFromSacService,
		},
		{
			name: "[new application flow] create application error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				accessPolicy := []string{"1", "2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{}, nil)
				sacClient.On("FindPoliciesByNames", accessPolicy).Return([]dto.PolicyDTO{}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:           appName,
					Site:           siteToUse,
					AccessPolicies: accessPolicy,
					Type:           model.HTTP,
				}
				sacClient.On("CreateApplication", &dto.ApplicationDTO{
					Name: appName,
					Type: model.HTTP,
				}).Return(&dto.ApplicationDTO{}, errorFromSacService)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    errorFromSacService,
		},
		{
			name: "[new application flow] bind application to site error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				accessPolicy := []string{"1", "2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{ID: "siteUUID"}, nil)
				sacClient.On("FindPoliciesByNames", accessPolicy).Return([]dto.PolicyDTO{}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:           appName,
					Site:           siteToUse,
					AccessPolicies: accessPolicy,
					Type:           model.HTTP,
				}
				sacClient.On("CreateApplication", &dto.ApplicationDTO{
					Name: appName,
					Type: model.HTTP,
				}).Return(&dto.ApplicationDTO{
					ID: "uuid",
				}, nil)
				sacClient.On("BindApplicationToSite", "uuid", "siteUUID").Return(errorFromSacService)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "uuid",
			},
			err: errorFromSacService,
		},
		{
			name: "[new application flow] bind application to site UpdatePolicies error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				accessPolicy := []string{"1", "2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{ID: "siteUUID"}, nil)
				sacClient.On("FindPoliciesByNames", accessPolicy).Return([]dto.PolicyDTO{
					{ID: "policy-uuid-1"},
					{ID: "policy-uuid-2"},
				}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:           appName,
					Site:           siteToUse,
					AccessPolicies: accessPolicy,
					Type:           model.HTTP,
				}
				sacClient.On("CreateApplication", &dto.ApplicationDTO{
					Name: appName,
					Type: model.HTTP,
				}).Return(&dto.ApplicationDTO{
					ID:   "uuid",
					Type: model.HTTP,
				}, nil)
				sacClient.On("BindApplicationToSite", "uuid", "siteUUID").Return(nil)
				sacClient.On("UpdatePolicies", "uuid", app.Type, []string{
					"policy-uuid-1",
					"policy-uuid-2",
				}).Return(errorFromSacService)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "uuid",
			},
			err: errorFromSacService,
		},
		{
			name: "[new application flow] success flow",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				accessPolicy := []string{"1", "2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{ID: "siteUUID"}, nil)
				sacClient.On("FindPoliciesByNames", accessPolicy).Return([]dto.PolicyDTO{
					{ID: "policy-uuid-1"},
					{ID: "policy-uuid-2"},
				}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:           appName,
					Site:           siteToUse,
					AccessPolicies: accessPolicy,
					Type:           model.HTTP,
				}
				sacClient.On("CreateApplication", &dto.ApplicationDTO{
					Name: appName,
					Type: model.HTTP,
				}).Return(&dto.ApplicationDTO{
					ID:   "uuid",
					Type: model.HTTP,
				}, nil)
				sacClient.On("BindApplicationToSite", "uuid", "siteUUID").Return(nil)
				sacClient.On("UpdatePolicies", "uuid", app.Type, []string{
					"policy-uuid-1",
					"policy-uuid-2",
				}).Return(nil)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "uuid",
			},
			err: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, input := test.setupFunc()
			output, err := s.Reconcile(context.Background(), input)
			assert.Equal(t, output, test.output)
			assert.ErrorIs(t, err, test.err)
		})
	}
}
