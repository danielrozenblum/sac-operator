package dto

import "github.com/google/uuid"

type SiteDTO struct {
	ID   *uuid.UUID `json:"id"`
	Name string     `json:"name"`
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
