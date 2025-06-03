package myError

import "errors"

type MyError struct {
	Message string
}

// 实现 Error 接口
func (e *MyError) Error() string {
	return e.Message
}

// NewMyError 创建自定义错误
func NewMyError(message string) error {
	return &MyError{Message: message}
}

// IsMyError 判断是否为自定义错误
func IsMyError(err error) bool {
	var myError *MyError
	ok := errors.As(err, &myError)
	return ok
}
