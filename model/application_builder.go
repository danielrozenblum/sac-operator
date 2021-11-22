package model

import (
	"github.com/google/uuid"
	"time"
)

type ApplicationBuilder struct {
	application *Application
}

func NewApplicationBuilder() *ApplicationBuilder {
	return &ApplicationBuilder{
		application: &Application{
			ID:         uuid.New(),
			Name:       "application-test",
			CreatedOn:  time.Now(),
			ModifiedOn: time.Now(),
		},
	}
}

func (a *ApplicationBuilder) WithID(id uuid.UUID) *ApplicationBuilder {
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
