package domain

type PaginationParams struct {
	Page  int
	Limit int
}

type PaginatedResponse[T any] struct {
	Items []T   `json:"items"`
	Total int   `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
