package response

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type Response struct {
	Status     string       `json:"status"`
	Data       any          `json:"data,omitempty"`
	Pagination *Pagination  `json:"pagination,omitempty"`
	Error      *ErrorDetail `json:"error,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}
