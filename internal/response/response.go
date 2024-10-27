package response

const OK_CODE = 200
const ERROR_CODE = 500

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code:    OK_CODE,
		Message: "success",
		Data:    data,
	}
}

func Fail(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

func FailWithMsg(message string) *Response {
	return &Response{
		Code:    ERROR_CODE,
		Message: message,
		Data:    nil,
	}
}
