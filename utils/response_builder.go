package utils

type ResponseEntity struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Ok(data interface{}) (resp *ResponseEntity) {
	return &ResponseEntity{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

func Error() (resp *ResponseEntity) {
	return &ResponseEntity{
		Code:    10001,
		Message: "server error",
		Data:    nil,
	}
}

func ErrorWithCodeMessage(code int, message string) (resp *ResponseEntity) {
	return &ResponseEntity{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

func ErrorWithMessage(message string) (resp *ResponseEntity) {
	return &ResponseEntity{
		Code:    10001,
		Message: message,
		Data:    nil,
	}
}
