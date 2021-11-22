package dto

type Site struct {
	Name string `json:"name"`
}

type Sites struct {
	First            bool   `json:"first"`
	Last             bool   `json:"last"`
	NumberOfElements int    `json:"numberOfElements"`
	Content          []Site `json:"content"`
	PageNumber       int    `json:"number"`
	PageSize         int    `json:"size"`
	TotalElements    int    `json:"totalElements"`
	TotalPages       int    `json:"totalPages"`
}
