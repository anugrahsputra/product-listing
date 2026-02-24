package dto

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type PaginatedResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}

type SuccessResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
