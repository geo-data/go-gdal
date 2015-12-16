package cpl

import (
	"strings"
)

type Error interface {
	error
	ErrorNum() int
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

func (e *cerr) ErrorNum() int {
	return e.n
}

func (e *cerr) ErrorType() int {
	return e.t
}

func NewError(cls, num int, msg string) Error {
	return &cerr{
		num,
		cls,
		msg,
	}
}

func SetError(err Error) {
	wrap_Error(err.ErrorType(), err.ErrorNum(), err.Error())
}

func LastError() (err Error) {
	t := GetLastErrorType()
	if t == 0 {
		return
	}

	err = NewError(
		t,
		GetLastErrorNo(),
		strings.TrimSpace(GetLastErrorMsg()),
	)
	return
}
