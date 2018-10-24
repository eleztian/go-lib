package errno

import "fmt"

type Errno struct {
	Code    int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

func (err *Errno) Add(msg string) *Errno {
	e := *err
	e.Message += " " + msg
	return &e
}

func (err *Errno) Addf(format string, args ...interface{}) *Errno {
	e := *err
	e.Message += " " + fmt.Sprintf(format, args...)
	return &e
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	return ErrInternalServer.Code, err.Error()
}
