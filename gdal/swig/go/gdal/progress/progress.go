package progress

/*
// Declare the GDALProgressFunc we're exporting from Go.
extern int GoGDALProgressFunc(double complete, char* message, void* data);
*/
import "C"
import "unsafe"

type ProgressFunc func(complete float64, message string, data interface{}) int

type args struct {
	pf   ProgressFunc
	data interface{}
}

func New(pf ProgressFunc, data interface{}) unsafe.Pointer {
	return unsafe.Pointer(&args{pf, data})
}

//export GoGDALProgressFunc
func GoGDALProgressFunc(complete C.double, message *C.char, data unsafe.Pointer) C.int {
	arg := (*args)(data)
	return C.int(arg.pf(
		float64(complete), C.GoString(message), arg.data,
	))
}

func GDALProgressFunc() unsafe.Pointer {
	return unsafe.Pointer(C.GoGDALProgressFunc)
}
