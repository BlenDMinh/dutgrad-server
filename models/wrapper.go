package models

type ResponseWrapper struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   *string     `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationResponse struct {
	ResponseWrapper
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
	Total    int64 `json:"total"`
}

func NewSuccessResponse(status int, message string, data interface{}) *ResponseWrapper {
	return &ResponseWrapper{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewPaginationResponse(status int, message string, data interface{}, page, pageSize, total int64) *PaginationResponse {
	return &PaginationResponse{
		ResponseWrapper: ResponseWrapper{
			Status:  status,
			Message: message,
			Data:    data,
		},
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
}

func NewErrorResponse(status int, message string, err *string) *ResponseWrapper {
	return &ResponseWrapper{
		Status:  status,
		Message: message,
		Error:   err,
	}
}
