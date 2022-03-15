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

func TestApplicationServiceImpl_Reconcile_DeleteApplication(t *testing.T) {
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

func TestApplicationServiceImpl_Reconcile_GetIDs(t *testing.T) {
	errorFromSacService := typederror.UnknownError
	tests := []struct {
		name      string
		setupFunc func() (ApplicationService, *model.Application)
		output    *ApplicationReconcileOutput
		err       error
	}{
		{
			name: "site does not exist",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{}, sac.ErrorNotFound)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:     appName,
					SiteName: siteToUse,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    typederror.UnrecoverableError,
		},
		{
			name: "find site by name error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{}, errorFromSacService)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:     appName,
					SiteName: siteToUse,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    errorFromSacService,
		},
		{
			name: "policy not exist",
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
					Name:                appName,
					SiteName:            siteToUse,
					AccessPoliciesNames: accessPolicy,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    typederror.UnrecoverableError,
		},
		{
			name: "find policy by name error",
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
					Name:                appName,
					SiteName:            siteToUse,
					AccessPoliciesNames: accessPolicy,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{},
			err:    errorFromSacService,
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

func TestApplicationServiceImpl_Reconcile_NewApplication(t *testing.T) {
	errorFromSacService := typederror.UnknownError
	tests := []struct {
		name      string
		setupFunc func() (ApplicationService, *model.Application)
		output    *ApplicationReconcileOutput
		err       error
	}{
		{
			name: "application already exist error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				accessPolicy := []string{"1", "2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{
					ID: "uuid",
				}, nil)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{ID: "siteUUID"}, nil)
				sacClient.On("FindPoliciesByNames", accessPolicy).Return([]dto.PolicyDTO{
					{ID: "policy-uuid-1"},
					{ID: "policy-uuid-2"},
				}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:                appName,
					SiteName:            siteToUse,
					AccessPoliciesNames: accessPolicy,
					Type:                model.HTTP,
				}
				sacClient.On("CreateApplication", &dto.ApplicationDTO{
					Name: appName,
					Type: model.HTTP,
				}).Return(&dto.ApplicationDTO{
					ID:   "uuid",
					Type: model.HTTP,
				}, errorFromSacService)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "",
			},
			err: typederror.UnrecoverableError,
		},
		{
			name: "find application error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				accessPolicy := []string{"1", "2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, errorFromSacService)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{ID: "siteUUID"}, nil)
				sacClient.On("FindPoliciesByNames", accessPolicy).Return([]dto.PolicyDTO{
					{ID: "policy-uuid-1"},
					{ID: "policy-uuid-2"},
				}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:                appName,
					SiteName:            siteToUse,
					AccessPoliciesNames: accessPolicy,
					Type:                model.HTTP,
				}
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "",
			},
			err: errorFromSacService,
		},
		{
			name: "create error",
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
					Name:                appName,
					SiteName:            siteToUse,
					AccessPoliciesNames: accessPolicy,
					Type:                model.HTTP,
				}
				sacClient.On("CreateApplication", &dto.ApplicationDTO{
					Name: appName,
					Type: model.HTTP,
				}).Return(&dto.ApplicationDTO{
					ID:   "uuid",
					Type: model.HTTP,
				}, errorFromSacService)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "",
			},
			err: errorFromSacService,
		},
		{
			name: "create success",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				accessPolicy := []string{"1", "2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindApplicationByName", appName).Return(&dto.ApplicationDTO{}, sac.ErrorNotFound)
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{ID: "site-uuid"}, nil)
				sacClient.On("FindPoliciesByNames", accessPolicy).Return([]dto.PolicyDTO{
					{ID: "policy-uuid-1"},
					{ID: "policy-uuid-2"},
				}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					Name:                appName,
					SiteName:            siteToUse,
					AccessPoliciesNames: accessPolicy,
					Type:                model.HTTP,
				}
				sacClient.On("CreateApplication", &dto.ApplicationDTO{
					Name: appName,
					Type: model.HTTP,
				}).Return(&dto.ApplicationDTO{
					ID:   "uuid",
					Type: model.HTTP,
				}, nil)
				sacClient.On("BindApplicationToSite", "uuid", "site-uuid").Return(nil)
				sacClient.On("UpdatePolicies", "uuid", model.ApplicationType("HTTP"), []string{"policy-uuid-1", "policy-uuid-2"}).Return(nil)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			output: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "uuid",
				SiteID:           "site-uuid",
				PoliciesIDs:      []string{"policy-uuid-1", "policy-uuid-2"},
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

func TestApplicationServiceImpl_Reconcile_UpdateApplication(t *testing.T) {
	errorFromSacService := typederror.UnknownError
	appID := "uuid"
	tests := []struct {
		name           string
		setupFunc      func() (ApplicationService, *model.Application)
		expectedOutput *ApplicationReconcileOutput
		err            error
	}{
		{
			name: "UpdateApplication error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				accessPolicy := []string{"access-policy-name-1", "access-policy-name-2"}
				activityPolicy := []string{"activity-policy-name-1", "activity-policy-name-2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{ID: "siteUUID"}, nil)
				sacClient.On("FindPoliciesByNames", append(accessPolicy, activityPolicy...)).Return([]dto.PolicyDTO{
					{ID: "policy-uuid-1"},
					{ID: "policy-uuid-2"},
				}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					ID:                    appID,
					Name:                  appName,
					Type:                  model.HTTP,
					SubType:               model.DefaultSubType,
					InternalAddress:       "internal",
					SiteName:              siteToUse,
					AccessPoliciesNames:   accessPolicy,
					ActivityPoliciesNames: activityPolicy,
					SiteId:                "",
					PoliciesIDS:           nil,
					ToDelete:              false,
				}
				sacClient.On("UpdateApplication", &dto.ApplicationDTO{
					ID:      "uuid",
					Name:    appName,
					Type:    model.HTTP,
					SubType: model.DefaultSubType,
					ConnectionSettings: dto.ConnectionSettingsDTO{
						InternalAddress: "internal",
						SubDomain:       "",
					},
					Icon:                  "",
					IsVisible:             false,
					IsNotificationEnabled: false,
					Enabled:               false,
				}).Return(&dto.ApplicationDTO{}, errorFromSacService)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			expectedOutput: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "uuid",
			},
			err: errorFromSacService,
		},
		{
			name: "site binding error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				siteID := "siteID"
				accessPolicy := []string{"access-policy-name-1", "access-policy-name-2"}
				activityPolicy := []string{"activity-policy-name-1", "activity-policy-name-2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{ID: siteID}, nil)
				sacClient.On("FindPoliciesByNames", append(accessPolicy, activityPolicy...)).Return([]dto.PolicyDTO{
					{ID: "policy-uuid-1"},
					{ID: "policy-uuid-2"},
				}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					ID:                    appID,
					Name:                  appName,
					Type:                  model.HTTP,
					SubType:               model.DefaultSubType,
					InternalAddress:       "internal",
					SiteName:              siteToUse,
					AccessPoliciesNames:   accessPolicy,
					ActivityPoliciesNames: activityPolicy,
					SiteId:                siteID,
					PoliciesIDS:           nil,
					ToDelete:              false,
				}
				sacClient.On("UpdateApplication", &dto.ApplicationDTO{
					ID:      "uuid",
					Name:    appName,
					Type:    model.HTTP,
					SubType: model.DefaultSubType,
					ConnectionSettings: dto.ConnectionSettingsDTO{
						InternalAddress: "internal",
						SubDomain:       "",
					},
					Icon:                  "",
					IsVisible:             false,
					IsNotificationEnabled: false,
					Enabled:               false,
				}).Return(&dto.ApplicationDTO{}, nil)
				sacClient.On("BindApplicationToSite", app.ID, app.SiteId).Return(errorFromSacService)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			expectedOutput: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "uuid",
			},
			err: errorFromSacService,
		},
		{
			name: "policy binding error",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				siteID := "siteID"
				accessPolicy := []string{"access-policy-name-1", "access-policy-name-2"}
				activityPolicy := []string{"activity-policy-name-1", "activity-policy-name-2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{ID: siteID}, nil)
				sacClient.On("FindPoliciesByNames", append(accessPolicy, activityPolicy...)).Return([]dto.PolicyDTO{
					{ID: "access-policy-uuid-1"},
					{ID: "access-policy-uuid-2"},
					{ID: "activity-policy-uuid-1"},
					{ID: "activity-policy-uuid-2"},
				}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					ID:                    appID,
					Name:                  appName,
					Type:                  model.ApplicationType("HTTP"),
					SubType:               model.DefaultSubType,
					InternalAddress:       "internal",
					SiteName:              siteToUse,
					AccessPoliciesNames:   accessPolicy,
					ActivityPoliciesNames: activityPolicy,
					SiteId:                siteID,
					PoliciesIDS:           nil,
					ToDelete:              false,
				}
				sacClient.On("UpdateApplication", &dto.ApplicationDTO{
					ID:      "uuid",
					Name:    appName,
					Type:    model.HTTP,
					SubType: model.DefaultSubType,
					ConnectionSettings: dto.ConnectionSettingsDTO{
						InternalAddress: "internal",
						SubDomain:       "",
					},
					Icon:                  "",
					IsVisible:             false,
					IsNotificationEnabled: false,
					Enabled:               false,
				}).Return(&dto.ApplicationDTO{}, nil)
				sacClient.On("BindApplicationToSite", app.ID, app.SiteId).Return(nil)
				sacClient.On("UpdatePolicies", app.ID, model.ApplicationType("HTTP"), []string{
					"access-policy-uuid-1",
					"access-policy-uuid-2",
					"activity-policy-uuid-1",
					"activity-policy-uuid-2",
				}).Return(errorFromSacService)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			expectedOutput: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "uuid",
			},
			err: errorFromSacService,
		},
		{
			name: "success flow",
			setupFunc: func() (ApplicationService, *model.Application) {
				appName := "test-application"
				siteToUse := "test-site"
				siteID := "siteID"
				accessPolicy := []string{"access-policy-name-1", "access-policy-name-2"}
				activityPolicy := []string{"activity-policy-name-1", "activity-policy-name-2"}
				sacClient := &sac.MockSecureAccessCloudClient{}
				sacClient.On("FindSiteByName", siteToUse).Return(&dto.SiteDTO{ID: siteID}, nil)
				sacClient.On("FindPoliciesByNames", append(accessPolicy, activityPolicy...)).Return([]dto.PolicyDTO{
					{ID: "access-policy-uuid-1"},
					{ID: "access-policy-uuid-2"},
					{ID: "activity-policy-uuid-1"},
					{ID: "activity-policy-uuid-2"},
				}, nil)
				testLog := ctrl.Log.WithName("test")
				app := &model.Application{
					ID:                    appID,
					Name:                  appName,
					Type:                  model.ApplicationType("HTTP"),
					SubType:               model.DefaultSubType,
					InternalAddress:       "internal",
					SiteName:              siteToUse,
					AccessPoliciesNames:   accessPolicy,
					ActivityPoliciesNames: activityPolicy,
					SiteId:                siteID,
					PoliciesIDS:           nil,
					ToDelete:              false,
				}
				sacClient.On("UpdateApplication", &dto.ApplicationDTO{
					ID:      "uuid",
					Name:    appName,
					Type:    model.HTTP,
					SubType: model.DefaultSubType,
					ConnectionSettings: dto.ConnectionSettingsDTO{
						InternalAddress: "internal",
						SubDomain:       "",
					},
					Icon:                  "",
					IsVisible:             false,
					IsNotificationEnabled: false,
					Enabled:               false,
				}).Return(&dto.ApplicationDTO{}, nil)
				sacClient.On("BindApplicationToSite", app.ID, app.SiteId).Return(nil)
				sacClient.On("UpdatePolicies", app.ID, model.ApplicationType("HTTP"), []string{
					"access-policy-uuid-1",
					"access-policy-uuid-2",
					"activity-policy-uuid-1",
					"activity-policy-uuid-2",
				}).Return(nil)
				return NewApplicationServiceImpl(sacClient, testLog), app
			},
			expectedOutput: &ApplicationReconcileOutput{
				Deleted:          false,
				SACApplicationID: "uuid",
				SiteID:           "siteID",
				PoliciesIDs: []string{
					"access-policy-uuid-1",
					"access-policy-uuid-2",
					"activity-policy-uuid-1",
					"activity-policy-uuid-2",
				},
			},
			err: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, input := test.setupFunc()
			output, err := s.Reconcile(context.Background(), input)
			assert.Equal(t, test.expectedOutput, output)
			assert.ErrorIs(t, err, test.err)
		})
	}
}
