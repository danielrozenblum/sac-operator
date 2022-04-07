package service

import (
	"context"
	"testing"

	"bitbucket.org/accezz-io/sac-operator/utils/typederror"

	connector_deployer "bitbucket.org/accezz-io/sac-operator/service/connector-deployer"

	"bitbucket.org/accezz-io/sac-operator/service/sac/dto"

	"github.com/stretchr/testify/mock"

	ctrl "sigs.k8s.io/controller-runtime"

	"bitbucket.org/accezz-io/sac-operator/model"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/accezz-io/sac-operator/service/sac"
)

func TestSiteServiceImpl_Reconcile(t *testing.T) {
	uncategorizedError := typederror.UnknownError
	tests := []struct {
		name      string
		setupFunc func() (SiteService, *model.Site)
		output    *SiteReconcileOutput
		err       error
	}{
		{
			name: "nil site flow",
			setupFunc: func() (SiteService, *model.Site) {
				sacClient := &sac.MockSecureAccessCloudClient{}
				testLog := ctrl.Log.WithName("test")
				return NewSiteServiceImpl(sacClient, nil, testLog), nil
			},
			output: &SiteReconcileOutput{},
			err:    typederror.UnrecoverableError,
		},
		{
			name: "delete site happy flow",
			setupFunc: func() (SiteService, *model.Site) {
				sacClient := &sac.MockSecureAccessCloudClient{}
				siteModel := &model.Site{
					ToDelete: true,
				}
				sacClient.On("DeleteSite", mock.AnythingOfType("string")).Return(nil)
				testLog := ctrl.Log.WithName("test")
				return NewSiteServiceImpl(sacClient, nil, testLog), siteModel
			},
			output: &SiteReconcileOutput{
				Deleted: true,
			},
			err: nil,
		},
		{
			name: "delete site failed flow",
			setupFunc: func() (SiteService, *model.Site) {
				sacClient := &sac.MockSecureAccessCloudClient{}
				siteModel := &model.Site{
					ToDelete: true,
				}
				sacClient.On("DeleteSite", mock.AnythingOfType("string")).Return(uncategorizedError)
				testLog := ctrl.Log.WithName("test")
				return NewSiteServiceImpl(sacClient, nil, testLog), siteModel
			},
			output: &SiteReconcileOutput{
				Deleted: false,
			},
			err: uncategorizedError,
		},
		{
			name: "create site success flow",
			setupFunc: func() (SiteService, *model.Site) {
				ctx := context.Background()
				sacClient := &sac.MockSecureAccessCloudClient{}
				deployer := &connector_deployer.MockConnectorDeployer{}
				siteModel := &model.Site{
					Name: "test",
				}
				siteDto := dto.FromSiteModel(siteModel)
				sacClient.On("CreateSite", siteDto).Return(&dto.SiteDTO{
					ID: "uuid",
				}, nil)
				deployer.On("GetConnectorsForSite", ctx, "test").Return([]connector_deployer.Connector{}, nil)
				testLog := ctrl.Log.WithName("test")
				return NewSiteServiceImpl(sacClient, deployer, testLog), siteModel
			},
			output: &SiteReconcileOutput{
				SACSiteID: "uuid",
			},
			err: nil,
		},
		{
			name: "create site err conflict flow",
			setupFunc: func() (SiteService, *model.Site) {
				sacClient := &sac.MockSecureAccessCloudClient{}
				siteModel := &model.Site{
					Name: "test",
				}
				siteDto := dto.FromSiteModel(siteModel)
				sacClient.On("CreateSite", siteDto).Return(&dto.SiteDTO{}, sac.ErrConflict)
				testLog := ctrl.Log.WithName("test")
				return NewSiteServiceImpl(sacClient, nil, testLog), siteModel
			},
			output: &SiteReconcileOutput{},
			err:    typederror.UnrecoverableError,
		},
		{
			name: "create site err flow",
			setupFunc: func() (SiteService, *model.Site) {
				sacClient := &sac.MockSecureAccessCloudClient{}
				siteModel := &model.Site{
					Name: "test",
				}
				siteDto := dto.FromSiteModel(siteModel)
				sacClient.On("CreateSite", siteDto).Return(&dto.SiteDTO{}, uncategorizedError)
				testLog := ctrl.Log.WithName("test")
				return NewSiteServiceImpl(sacClient, nil, testLog), siteModel
			},
			output: &SiteReconcileOutput{},
			err:    uncategorizedError,
		},
		{
			name: "failed get connector flow",
			setupFunc: func() (SiteService, *model.Site) {
				ctx := context.Background()
				deployer := &connector_deployer.MockConnectorDeployer{}
				siteModel := &model.Site{
					Name:      "test",
					SACSiteID: "uuid",
				}
				connectorList := []connector_deployer.Connector{}
				deployer.On("GetConnectorsForSite", ctx, "test").Return(connectorList, uncategorizedError)
				testLog := ctrl.Log.WithName("test")
				return NewSiteServiceImpl(nil, deployer, testLog), siteModel
			},
			output: &SiteReconcileOutput{
				SACSiteID: "uuid",
			},
			err: uncategorizedError,
		},
		{
			name: "failed get connector flow",
			setupFunc: func() (SiteService, *model.Site) {
				ctx := context.Background()
				deployer := &connector_deployer.MockConnectorDeployer{}
				siteModel := &model.Site{
					Name:      "test",
					SACSiteID: "uuid",
				}
				connectorList := []connector_deployer.Connector{}
				deployer.On("GetConnectorsForSite", ctx, "test").Return(connectorList, uncategorizedError)
				testLog := ctrl.Log.WithName("test")
				return NewSiteServiceImpl(nil, deployer, testLog), siteModel
			},
			output: &SiteReconcileOutput{
				SACSiteID: "uuid",
			},
			err: uncategorizedError,
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

func TestSiteServiceImpl_connectorDeploymentArgsFromCommand(t *testing.T) {
	service := NewSiteServiceImpl(nil, nil, nil)

	command := &dto.ConnectorDeploymentCommand{DeploymentCommands: "test-ps8l:\n    image: luminate/connector:2.10.1\n    container_name: test-ps8l\n    restart: on-failure\n    ulimits:\n        nofile: 2048\n    log_opt:\n        max-size: \"50m\"\n        max-file: \"10\"\n    environment:\n     - ENDPOINT_URL=symchatbotdemo.luminatesite.com\n     - TENANT_IDENTIFIER=1667d9c1d754419a9c456f47b9c1df28_symchatbotdemo\n     - HTTPS_SKIP_CERT_VERIFY=true\n     - OTP=7b86d374-23c4-4b32-a156-f80ec5909f5b\n     "}

	args, err := service.connectorDeploymentArgsFromCommand(command)
	assert.Nil(t, err)
	assert.Equal(t, "luminate/connector:2.10.1", args.Image)
	assert.Equal(t, "test-ps8l", args.ContainerName)
	assert.Equal(t, args.EnvironmentVars["ENDPOINT_URL"], "symchatbotdemo.luminatesite.com")
}
