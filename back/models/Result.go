package models

type Result struct {
	Token   string
	Status  int // 0失败，1成功
	Message string
	Data    interface{}
}

func Fail(token, message string, data interface{}) (result Result) {
	result.Token = token
	result.Status = 0
	result.Message = message
	result.Data = data
	return
}

func Success(token, message string, data interface{}) (result Result) {
	result.Token = token
	result.Status = 1
	result.Message = message
	result.Data = data
	return
}
