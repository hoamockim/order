package errors

type AppError struct {
	meta   *ErrorMeta   `json:"meta"`
	detail *ErrorDetail `json:"detail"`
}

type ErrorMeta struct {
	HttpCode int    `json:"http_code"`
	Code     string `json:"code"`
	Key      string `json:"key"`
}

type ErrorDetail struct {
	messages []string `json:"message"`
	service  string   `json:"service"`
}

func (detail *ErrorDetail) SetMessageError(message string) {
	detail.messages = append(detail.messages, message)
}

func getAppErr(code string, key string) AppError {
	return AppError{
		meta: &ErrorMeta{
			Code: code,
			Key:  key,
		},
	}
}

func New(errMeta *ErrorMeta, serviceName string) *AppError {
	return &AppError{
		meta: errMeta,
		detail: &ErrorDetail{
			service: serviceName,
		},
	}
}
