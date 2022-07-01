package errno

import (
	"encoding/json"
	"fmt"
	"net/http"

	errors "github.com/micro/go-micro/errors"
)

type ErrDetail struct {
	Errno
	Cause string `json:"cause"`
}

func NotFoundErr(errno *Errno, cause string) error {
	detail := &ErrDetail{
		*errno,
		cause,
	}

	detailStr, _ := json.Marshal(detail)

	return &errors.Error{
		Id:     "serverName",
		Code:   404,
		Detail: string(detailStr),
		Status: http.StatusText(404),
	}
}

func BadRequestErr(errno *Errno, cause string) error {
	detail := &ErrDetail{
		*errno,
		cause,
	}

	detailStr, _ := json.Marshal(detail)

	return &errors.Error{
		Id:     "serverName",
		Code:   400,
		Detail: string(detailStr),
		Status: http.StatusText(400),
	}
}

func ServerErr(errno *Errno, cause string) error {
	detail := &ErrDetail{
		*errno,
		cause,
	}

	detailStr, _ := json.Marshal(detail)

	return &errors.Error{
		Id:     "serverName",
		Code:   500,
		Detail: string(detailStr),
		Status: http.StatusText(500),
	}
}

func ParseDetail(detail string) (*ErrDetail, error) {
	d := &ErrDetail{}
	err := json.Unmarshal([]byte(detail), d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (err Errno) Error() string {
	return err.Message
}

// // Err represents an error
type Err struct {
	Code    int
	Message string
	Err     error
}

func New(errno *Errno, err error) *Err {
	return &Err{Code: errno.Code, Message: errno.Message, Err: err}
}

func (err *Err) Add(message string) error {
	err.Message += " " + message
	return err
}

func (err *Err) Addf(format string, args ...interface{}) error {
	err.Message += " " + fmt.Sprintf(format, args...)
	return err
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

func IsErrUserNotExisted(err error) bool {
	code, _ := DecodeErr(err)
	return code == ErrUserNotExisted.Code
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	return InternalServerError.Code, err.Error()
}
