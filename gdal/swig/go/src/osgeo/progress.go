package gdal

import "C"
import "unsafe"

type ProgressFunc func(complete float64, message string, progressArg interface{}) int

type goGDALProgressFuncProxyArgs struct {
	progresssFunc ProgressFunc
	data          interface{}
}

//export goGDALProgressFuncProxyA
func goGDALProgressFuncProxyA(complete C.double, message *C.char, data unsafe.Pointer) C.int {
	arg := (*goGDALProgressFuncProxyArgs)(data)
	return C.int(arg.progresssFunc(
		float64(complete), C.GoString(message), arg.data,
	))
}
