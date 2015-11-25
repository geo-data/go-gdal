package gdal

import (
	"fmt"
	"strings"
)

type Err interface {
	error
	ErrorNo() int
	ErrorType() int
}

type gerr struct {
	n int    // error number
	t int    // error type
	m string // error message
}

func (e *gerr) Error() string {
	return e.m
}

func (e *gerr) ErrorNo() int {
	return e.n
}

func (e *gerr) ErrorType() int {
	return e.t
}

func lastError() (err Err) {
	t := GetLastErrorType()
	if t == 0 {
		return
	}

	err = &gerr{
		GetLastErrorNo(),
		t,
		strings.TrimSpace(GetLastErrorMsg()),
	}
	ErrorReset()
	return
}

func SetErrorHandler(callbackName string) (err error) {
	if wrap_SetErrorHandler(callbackName) != 0 {
		err = fmt.Errorf("Non existent error handler: %s", callbackName)
	}
	return
}
