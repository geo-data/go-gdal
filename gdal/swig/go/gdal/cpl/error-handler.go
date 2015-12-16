package cpl

/*

// Declare CPLGetErrorHandlerUserData()
#include "cpl_error.h"

// Declare the CPLErrorHandler we're exporting from Go.
extern void GoCPLErrorHandler(int errcls, int errnum, char* message);
*/
import "C"
import "unsafe"

type ErrorHandler func(err Error, data interface{})

type handlerargs struct {
	handler ErrorHandler
	data    interface{}
}

func NewErrorHandler(handler ErrorHandler, data interface{}) unsafe.Pointer {
	return unsafe.Pointer(&handlerargs{handler, data})
}

//export GoCPLErrorHandler
func GoCPLErrorHandler(errcls C.int, errnum C.int, message *C.char) {
	arg := (*handlerargs)(C.CPLGetErrorHandlerUserData())
	arg.handler(
		NewError(int(errcls), int(errnum), C.GoString(message)),
		arg.data,
	)
	return
}

func CPLErrorHandler() unsafe.Pointer {
	return unsafe.Pointer(C.GoCPLErrorHandler)
}
