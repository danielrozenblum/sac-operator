package model

import (
	"github.com/google/uuid"
)

type ApplicationBuilder struct {
	application *Application
}

func NewApplicationBuilder() *ApplicationBuilder {
	applicationId := uuid.New().String()

	return &ApplicationBuilder{
		application: &Application{
			ID:                    applicationId,
			Name:                  "application-test",
			Type:                  HTTP,
			SubType:               DefaultSubType,
			InternalAddress:       "http://1.1.1.1",
			SiteName:              "site-test-1",
			AccessPoliciesNames:   []string{"access-policy-1", "access-policy-2"},
			ActivityPoliciesNames: []string{},
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
