package models

type Response struct {
	Success	bool		`json:"succes"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Success:	true,
		Data:		data,
	}
}

func NewErrorResponse(code int, message string) *Response {
	return &Response{
		Success: false,
		Error:	&ErrorInfo{
			Code:	code,
			Message: message,
		},
	}
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination,omitempty"`
	Error      *ErrorInfo  `json:"error,omitempty"`
}

func NewPaginatedResponse(data interface{}, page, pageSize, total int) *PaginatedResponse {
	return &PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}
}