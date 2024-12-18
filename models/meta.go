package models

type (
	Meta struct {
		TotalRecords int `json:"total"`
		Page         int `json:"page"`
		TotalPages   int `json:"pages"`
		Limit        int `json:"limit"`
	}
)
