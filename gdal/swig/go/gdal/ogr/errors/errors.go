package errors

import "fmt"

type Error interface {
	error
	ErrorNo() int
}

type ogrerr struct {
	n int    // error number
	m string // error message
}

func (e *ogrerr) Error() string {
	return e.m
}

func (e *ogrerr) ErrorNo() int {
	return e.n
}

func New(num int, msg string) (err Error) {
	return &ogrerr{
		num,
		msg,
	}
}

func Errorf(num int, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return New(num, msg)
}
