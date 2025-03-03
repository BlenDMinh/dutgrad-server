package models

type ResponseWrapper struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   *string     `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(status int, message string, data interface{}) *ResponseWrapper {
	return &ResponseWrapper{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(status int, message string, err *string) *ResponseWrapper {
	return &ResponseWrapper{
		Status:  status,
		Message: message,
		Error:   err,
	}
}
