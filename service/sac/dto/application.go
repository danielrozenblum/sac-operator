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
	Icon                  *string                  `json:"icon"`
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
	return &ApplicationDTO{
		ID:      application.ID,
		Name:    application.Name,
		Type:    application.Type,
		SubType: application.SubType,
		ConnectionSettings: ConnectionSettingsDTO{
			InternalAddress: application.InternalAddress,
		},
	}
}

func ToApplicationModel(dto *ApplicationDTO, siteID string) *model.Application {
	return &model.Application{
		ID:               dto.ID,
		Name:             dto.Name,
		Type:             dto.Type,
		SubType:          dto.SubType,
		InternalAddress:  dto.ConnectionSettings.InternalAddress,
		Site:             siteID,
		AccessPolicies:   nil,
		ActivityPolicies: nil,
	}
}

func MergeApplication(existingApplication *ApplicationDTO, updatedApplication *ApplicationDTO) *ApplicationDTO {
	mergedApplication := *existingApplication

	mergedApplication.Name = updatedApplication.Name
	mergedApplication.Type = updatedApplication.Type
	mergedApplication.SubType = updatedApplication.SubType
	mergedApplication.ConnectionSettings.InternalAddress = updatedApplication.ConnectionSettings.InternalAddress

	return &mergedApplication
}
