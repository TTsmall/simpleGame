package errorx

import "fmt"

type ErrorItem struct {
	Code    int32
	Message string
}

func NewErrorItem(format string, args ...interface{}) *ErrorItem {
	return &ErrorItem{-999, fmt.Sprintf(format, args...)}
}

func Create(code int32, message string) *ErrorItem {
	return &ErrorItem{Code: code, Message: message}
}

func (this *ErrorItem) Error() string {
	return this.Message
}

func (this *ErrorItem) CloneWithMsg(format string, args ...interface{}) *ErrorItem {
	return &ErrorItem{
		Code:    this.Code,
		Message: fmt.Sprintf(format, args...),
	}
}

func (this *ErrorItem) CloneWithArgs(args ...interface{}) *ErrorItem {
	return &ErrorItem{
		Code:    this.Code,
		Message: fmt.Sprintf(this.Message, args...),
	}
}
