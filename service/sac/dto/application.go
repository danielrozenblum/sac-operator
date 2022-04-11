package dto

import (
	"bitbucket.org/accezz-io/sac-operator/model"
	"github.com/jinzhu/copier"
)

type ApplicationDTO struct {
	ID                    string                   `json:"id,omitempty"`
	Name                  string                   `json:"name,omitempty"`
	Type                  model.ApplicationType    `json:"type,omitempty"`
	SubType               model.ApplicationSubType `json:"subType,omitempty"`
	IconUrl               string                   `json:"iconUrl,omitempty"`
	IsVisible             bool                     `json:"isVisible"`
	IsNotificationEnabled bool                     `json:"isNotificationEnabled"`
	Enabled               bool                     `json:"enabled"`

	ConnectionSettings               ConnectionSettingsDTO             `json:"connectionSettings"`
	HttpLinkTranslationSettings      *HttpLinkTranslationSettingsDTO   `json:"linkTranslationSettings,omitempty"`
	HttpRequestCustomizationSettings *HttpRequestCustomizationSettings `json:"requestCustomizationSettings,omitempty"`
}

type ConnectionSettingsDTO struct {
	InternalAddress       string `json:"internalAddress"`
	Subdomain             string `json:"subdomain"`
	LuminateAddress       string `json:"luminateAddress,omitempty"`
	ExternalAddress       string `json:"externalAddress,omitempty"`
	CustomExternalAddress string `json:"customExternalAddress,omitempty"`
	CustomRootPath        string `json:"customRootPath,omitempty"`
	HealthUrl             string `json:"healthUrl,omitempty"`
	HealthMethod          string `json:"healthMethod,omitempty"`
	CustomSSLCertificate  string `json:"customSSLCertificate,omitempty"`
	WildcardPrivateKey    string `json:"wildcardPrivateKey,omitempty"`
}

type HttpLinkTranslationSettingsDTO struct {
	IsDefaultContentRewriteRulesEnabled bool     `json:"isDefaultContentRewriteRulesEnabled"`
	IsDefaultHeaderRewriteRulesEnabled  bool     `json:"isDefaultHeaderRewriteRulesEnabled"`
	UseExternalAddressForHostAndSni     bool     `json:"useExternalAddressForHostAndSni"`
	LinkedApplications                  []string `json:"linkedApplications"`
}

type HttpRequestCustomizationSettings struct {
	HeaderCustomization map[string]string `json:"headerCustomization"`
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

func FromApplicationModel(application *model.Application) (*ApplicationDTO, error) {
	dto := &ApplicationDTO{
		ID:                    application.ID,
		Name:                  application.Name,
		Type:                  application.Type,
		SubType:               application.SubType,
		IsVisible:             application.IsVisible,
		IsNotificationEnabled: application.IsNotificationEnabled,
		Enabled:               application.Enabled,
	}

	dto.ConnectionSettings = ConnectionSettingsDTO{}
	if application.ConnectionSettings != nil {
		err := copier.Copy(&dto.ConnectionSettings, application.ConnectionSettings)
		if err != nil {
			return nil, err
		}
	}
	if application.HttpLinkTranslationSettings != nil {
		dto.HttpLinkTranslationSettings = &HttpLinkTranslationSettingsDTO{}
		err := copier.Copy(dto.HttpLinkTranslationSettings, application.HttpLinkTranslationSettings)
		if err != nil {
			return nil, err
		}
	}
	if application.HttpRequestCustomizationSettings != nil {
		dto.HttpRequestCustomizationSettings = &HttpRequestCustomizationSettings{}
		err := copier.Copy(dto.HttpRequestCustomizationSettings, application.HttpRequestCustomizationSettings)
		if err != nil {
			return nil, err
		}
	}

	return dto, nil
}

type MergeOptions struct {
	LinkTranslationSettings      bool
	RequestCustomizationSettings bool
	Icon                         bool
	ConnectionSettings           bool
}

func MergeApplication(existingApplication *ApplicationDTO, updatedApplication *ApplicationDTO, options MergeOptions) *ApplicationDTO {
	mergedApplication := *existingApplication

	mergedApplication.Name = updatedApplication.Name
	mergedApplication.Type = updatedApplication.Type
	mergedApplication.SubType = updatedApplication.SubType
	mergedApplication.IconUrl = updatedApplication.IconUrl
	mergedApplication.ConnectionSettings.InternalAddress = updatedApplication.ConnectionSettings.InternalAddress

	return &mergedApplication
}
