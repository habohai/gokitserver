package util

type MyError struct {
	Code    int
	Message string
}

func NewMyError(code int, msg string) error {
	return &MyError{
		Code:    code,
		Message: msg,
	}
}

func (m *MyError) Error() string {
	return m.Message
}
