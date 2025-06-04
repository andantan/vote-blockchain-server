package werror

import "fmt"

type WrappedError struct {
	Code    string
	Message string
	Err     error
}

func NewWrappedError(code, message string, err error) *WrappedError {
	return &WrappedError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *WrappedError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}

	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *WrappedError) Unwrap() error {
	return e.Err
}
