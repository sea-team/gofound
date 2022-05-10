package web

type Result struct {
	State bool `json:"state"`

	Message string `json:"message,omitempty"`

	Data interface{} `json:"data,omitempty"`
}

func NewResult(state bool, message string, data interface{}) *Result {
	return &Result{
		State:   state,
		Message: message,
		Data:    data,
	}
}

// Success 返回成功
func Success(data interface{}) *Result {
	return NewResult(true, "success", data)
}

// Error 返回错误
func Error(message string) *Result {
	return NewResult(false, message, nil)
}
