package myError

import "errors"

type myError struct {
	Message string
}

// 实现 Error 接口
func (e *myError) Error() string {
	return e.Message
}

// NewMyError 创建自定义错误
func NewMyError(message string) error {
	return &myError{Message: message}
}

// IsMyError 判断是否为自定义错误
func IsMyError(err error) bool {
	var myErr *myError
	ok := errors.As(err, &myErr)
	return ok
}
