package cpl

import (
	"fmt"
	"strings"
)

type Error interface {
	error
	ErrorNo() int
	ErrorType() int
}

type cerr struct {
	n int    // error number
	t int    // error type
	m string // error message
}

func (e *cerr) Error() string {
	return e.m
}

func (e *cerr) ErrorNo() int {
	return e.n
}

func (e *cerr) ErrorType() int {
	return e.t
}

func LastError() (err Error) {
	t := GetLastErrorType()
	if t == 0 {
		return
	}

	err = &cerr{
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
