package dto

import (
	"bitbucket.org/accezz-io/sac-operator/model"
	"github.com/google/uuid"
)

type ApplicationDTO struct {
	ID                    *uuid.UUID               `json:"id"`
	Name                  string                   `json:"name"`
	Type                  model.ApplicationType    `json:"type"`
	SubType               model.ApplicationSubType `json:"subType"`
	ConnectionSettings    ConnectionSettingsDTO    `json:"connectionSettings"`
	Icon                  string                   `json:"icon"`
	IsVisible             bool                     `json:"isVisible"`
	IsNotificationEnabled bool                     `json:"isNotificationEnabled"`
	Enabled               bool                     `json:"enabled"`
}

type ConnectionSettingsDTO struct {
	InternalAddress string `json:"internalAddress"`
	SubDomain       string `json:"subDomain"`
}

type ApplicationPageDTO struct {
	First            bool             `json:"first"`
	Last             bool             `json:"last"`
	NumberOfElements int              `json:"numberOfElements"`
	Content          []ApplicationDTO `json:"content"`
	PageNumber       int              `json:"number"`
	PageSize         int              `json:"size"`
	TotalElements    int              `json:"totalElements"`
	TotalPages       int              `json:"totalPages"`
}

func FromApplicationModel(application *model.Application) *ApplicationDTO {
	// TODO: implement
	return nil
}

func ToApplicationModel(dto *ApplicationDTO) *model.Application {
	// TODO: implement
	return nil
}
