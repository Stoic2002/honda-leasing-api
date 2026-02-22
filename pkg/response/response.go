package response

type StandardResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page    int   `json:"page"`
	Limit   int   `json:"limit"`
	HasMore bool  `json:"has_more"`
	Total   int64 `json:"total,omitempty"`
}

func Success(code int, message string, data interface{}) StandardResponse {
	return StandardResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func SuccessPaginated(code int, message string, data interface{}, meta PaginationMeta) PaginatedResponse {
	return PaginatedResponse{
		Code:    code,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func Error(code int, message string) StandardResponse {
	return StandardResponse{
		Code:    code,
		Message: message,
	}
}
