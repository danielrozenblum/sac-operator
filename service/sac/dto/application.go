package dto

type Application struct {
	Name string `json:"name"`
}

type Applications struct {
	First            bool          `json:"first"`
	Last             bool          `json:"last"`
	NumberOfElements int           `json:"numberOfElements"`
	Content          []Application `json:"content"`
	PageNumber       int           `json:"number"`
	PageSize         int           `json:"size"`
	TotalElements    int           `json:"totalElements"`
	TotalPages       int           `json:"totalPages"`
}
