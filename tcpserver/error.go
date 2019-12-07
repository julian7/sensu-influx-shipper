package tcpserver

import "fmt"

type InAction interface {
	Action() string
	Value() interface{}
}

type Error struct {
	action string
	err    error
	value  interface{}
}

func (e *Error) Action() string {
	return e.action
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Value() interface{} {
	return e.value
}

func NewError(action string, err error) error {
	if err == nil {
		return nil
	}

	return &Error{action: action, err: err}
}

func NewErrorWithValue(action string, value interface{}, err error) error {
	if err == nil {
		return nil
	}

	return &Error{action: action, err: err, value: value}
}

func Errorf(action, format string, a ...interface{}) error {
	return NewError(action, fmt.Errorf(format, a...))
}
