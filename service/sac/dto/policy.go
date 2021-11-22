package dto

import "time"

type Policies struct {
	First            bool     `json:"first"`
	Last             bool     `json:"last"`
	NumberOfElements int      `json:"numberOfElements"`
	Number           int      `json:"number"`
	Content          []Policy `json:"content"`
	Size             int      `json:"size"`
	TotalElements    int      `json:"totalElements"`
	TotalPages       int      `json:"totalPages"`
}

type Policy struct {
	AllApplications      bool `json:"allApplications"`
	AllDirectoryEntities bool `json:"allDirectoryEntities"`
	Applications         []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"applications"`
	Containers        []interface{}     `json:"containers"`
	CreatedAt         time.Time         `json:"createdAt"`
	DirectoryEntities []DirectoryEntity `json:"directoryEntities"`
	Enabled           bool              `json:"enabled"`
	FilterConditions  []interface{}     `json:"filterConditions"`
	ID                string            `json:"id"`
	IsDefault         bool              `json:"isDefault"`
	ModifiedOn        time.Time         `json:"modifiedOn"`
	Name              string            `json:"name"`
	Static            bool              `json:"static"`
	TargetProtocol    string            `json:"targetProtocol"`
	Type              string            `json:"type"`
	Validators        struct {
	} `json:"validators"`
}
