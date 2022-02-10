package dto

import (
	"time"

	"bitbucket.org/accezz-io/sac-operator/model"
	"github.com/google/uuid"
)

type SiteDTO struct {
	ID         *uuid.UUID  `json:"id"`
	Name       string      `json:"name"`
	Connectors []Connector `json:"connectors"`
}

type SitePageDTO struct {
	First            bool      `json:"first"`
	Last             bool      `json:"last"`
	NumberOfElements int       `json:"numberOfElements"`
	Content          []SiteDTO `json:"content"`
	PageNumber       int       `json:"number"`
	PageSize         int       `json:"size"`
	TotalElements    int       `json:"totalElements"`
	TotalPages       int       `json:"totalPages"`
}

func FromSiteModel(site *model.Site) *SiteDTO {
	return &SiteDTO{
		ID:   site.ID,
		Name: site.Name,
	}
}

type Connector struct {
	ID                             string     `json:"id,omitempty"`
	Name                           string     `json:"name,omitempty"`
	Otp                            string     `json:"otp,omitempty"`
	DateCreated                    *time.Time `json:"date_created,omitempty"`
	DateRegistered                 *time.Time `json:"date_registered,omitempty"`
	DateOtpExpire                  *time.Time `json:"date_otp_expire,omitempty"`
	SendLogs                       bool       `json:"send_logs,omitempty"`
	Enabled                        bool       `json:"enabled,omitempty"`
	ConnectorStatus                string     `json:"connector_status,omitempty"`
	UpdateStatus                   string     `json:"update_status,omitempty"`
	UpdateStatusInfo               string     `json:"update_status_info,omitempty"`
	InternalIP                     string     `json:"internal_ip,omitempty"`
	ExternalIP                     string     `json:"external_ip,omitempty"`
	Hostname                       string     `json:"hostname,omitempty"`
	GeoLocation                    string     `json:"geo_location,omitempty"`
	DeploymentType                 string     `json:"deployment_type,omitempty"`
	KubernetesPersistentVolumeName string     `json:"kubernetes_persistent_volume_name,omitempty"`
	Version                        string     `json:"version"`
}

type ConnectorPageDTO struct {
	First            bool        `json:"first"`
	Last             bool        `json:"last"`
	NumberOfElements int         `json:"numberOfElements"`
	Content          []Connector `json:"content"`
	PageNumber       int         `json:"number"`
	PageSize         int         `json:"size"`
	TotalElements    int         `json:"totalElements"`
	TotalPages       int         `json:"totalPages"`
}
