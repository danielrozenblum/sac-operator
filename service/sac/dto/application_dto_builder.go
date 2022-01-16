package dto

import (
	"bitbucket.org/accezz-io/sac-operator/model"
	"github.com/google/uuid"
)

type ApplicationDTOBuilder struct {
	dto *ApplicationDTO
}

func NewApplicationDTOBuilder() *ApplicationDTOBuilder {
	applicationId := uuid.New()

	return &ApplicationDTOBuilder{dto: &ApplicationDTO{
		ID:                    &applicationId,
		Name:                  "test",
		Type:                  model.DefaultType,
		SubType:               model.DefaultSubType,
		ConnectionSettings:    ConnectionSettingsDTO{},
		Icon:                  nil,
		IsVisible:             true,
		IsNotificationEnabled: true,
		Enabled:               true,
	}}
}

func (a *ApplicationDTOBuilder) WithID(id *uuid.UUID) *ApplicationDTOBuilder {
	a.dto.ID = id
	return a
}

func (a *ApplicationDTOBuilder) WithName(name string) *ApplicationDTOBuilder {
	a.dto.Name = name
	return a
}

func (a *ApplicationDTOBuilder) WithIsVisible(isVisible bool) *ApplicationDTOBuilder {
	a.dto.IsVisible = isVisible
	return a
}

func (a *ApplicationDTOBuilder) Build() *ApplicationDTO {
	return a.dto
}
