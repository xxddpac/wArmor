package errors

import "fmt"

//func NewBusinessError(defaultMessage string, code int) BusinessError {
//	return BusinessError{
//		DefaultMessage: defaultMessage,
//		Code:           code,
//	}
//}

type BusinessError struct {
	DefaultMessage string `Json:"default message"`
	message        string
	Code           int `Json:"business code"`
	rawError       error
}

func (b BusinessError) AttachError(rawError error) BusinessError {
	b.rawError = rawError
	return b
}

func (b BusinessError) AttachMessage(message interface{}) BusinessError {
	b.message = fmt.Sprint(message)
	return b
}

func (b BusinessError) AppendMessage(message interface{}) BusinessError {
	b.message = b.message + fmt.Sprint(message)
	return b
}

func (b BusinessError) Error() string {
	errorMessage := b.DefaultMessage
	if nil != b.rawError {
		errorMessage += ". raw error: " + b.rawError.Error()
	}
	if "" != b.message {
		errorMessage += ". message: " + b.message
	}

	return errorMessage
}

func (b BusinessError) Bytes() []byte {
	return []byte(b.Error())
}
