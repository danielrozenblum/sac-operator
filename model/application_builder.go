package model

import (
	"github.com/google/uuid"
)

type ApplicationBuilder struct {
	application *Application
}

func NewApplicationBuilder() *ApplicationBuilder {
	applicationId := uuid.New().String()

	common := CommonApplicationParams{
		Name:                  "application-test",
		SiteName:              "",
		IsVisible:             false,
		IsNotificationEnabled: false,
		Enabled:               false,
		AccessPoliciesNames:   []string{"access-policy-1", "access-policy-2"},
		ActivityPoliciesNames: []string{},
	}

	return &ApplicationBuilder{
		application: &Application{
			ID:                      applicationId,
			Type:                    HTTP,
			SubType:                 DefaultSubType,
			CommonApplicationParams: common,
		},
	}
}

func (a *ApplicationBuilder) WithID(id string) *ApplicationBuilder {
	a.application.ID = id
	return a
}

func (a *ApplicationBuilder) WithName(name string) *ApplicationBuilder {
	a.application.Name = name
	return a
}

func (a *ApplicationBuilder) Build() *Application {
	return a.application
}
