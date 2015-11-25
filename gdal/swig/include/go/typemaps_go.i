%include "typemaps.i"

%typemap(imtype) IF_FALSE_RETURN_NONE "int"

%typemap(gotype) (double argin[ANY]), (double argout[ANY]) %{[]float64%}
%typemap(imtype) (double argin[ANY]), (double argout[ANY]) "*C.double"
%typemap(goin) (double argin[ANY]), (double argout[ANY]) %{
  $result = (*C.double)(unsafe.Pointer(&$input[0]))
%}

/*%typemap(goin) (double argout[ANY]) %{
    $result = (*C.double)(unsafe.Pointer(&$input[0]))
    %}*/
/*%typemap(goout) (double argout[ANY]) %{
  if $input != nil {
      $input_p := $input
	$input_q := uintptr(unsafe.Pointer($input_p))
	for {
		$input_p = (**C.double)(unsafe.Pointer($input_q))
		if *$input_p == nil {
			break
		}
		$result = append($result, float64(*$input_p))
		$input_q += unsafe.Sizeof($input_q)
	}
    }
    }
%}*/

%typemap(gotype) char **options %{[]string%}
%typemap(goin) char **options %{
  	$input_l := len($input)
    $result_i := make([]*C.char, $input_l+1)
	for i := 0; i < $input_l; i++ {
		$result_i[i] = C.CString($input[i])
		defer C.free(unsafe.Pointer($result_i[i]))
	}
	$result_i[$input_l] = (*C.char)(unsafe.Pointer(nil))
    $result = &$result_i[0]
%}
%typemap(imtype) char **options "**C.char"

%typemap(gotype) char ** %{[]string%}
%typemap(goout) char **options %{
  if $input != nil {
    $input_p := $input
	$input_q := uintptr(unsafe.Pointer($input_p))
	for {
		$input_p = (**C.char)(unsafe.Pointer($input_q))
		if *$input_p == nil {
			break
		}
		$result = append($result, C.GoString(*$input_p))
		$input_q += unsafe.Sizeof($input_q)
	}
    }
%}

%typemap(gotype) char **CSL %{[]string%}
%typemap(imtype) char **CSL "**C.char"
%typemap(goout) char **CSL %{
  if $input != nil {
    defer C.CSLDestroy($input)
    $input_p := $input
	$input_q := uintptr(unsafe.Pointer($input_p))
	for {
		$input_p = (**C.char)(unsafe.Pointer($input_q))
		if *$input_p == nil {
			break
		}
		$result = append($result, C.GoString(*$input_p))
		$input_q += unsafe.Sizeof($input_q)
	}
    }
%}

%typemap(gotype) ( void* callback_data = NULL ) "interface{}"
%typemap(imtype) ( void* callback_data = NULL) "unsafe.Pointer"
%typemap(gotype) ( GDALProgressFunc callback = NULL ) "ProgressFunc"
%typemap(imtype) ( GDALProgressFunc callback = NULL ) "C.GDALProgressFunc"
%typemap(goin) ( GDALProgressFunc callback = NULL ) %{
    var progress ProgressFunc
	if $input != nil {
        progress = $input
   		$result = C.GDALProgressFunc(C.goGDALProgressFuncProxyA)
    }
%}
%typemap(goin) ( void* callback_data = NULL ) %{
    if progress != nil {
        $result = unsafe.Pointer(&goGDALProgressFuncProxyArgs{
            progress, $input,
        })
    }
%}
%typemap(in) ( GDALProgressFunc callback = NULL) 
{
  $1 = $input;
}

%define CHECK_NULL()
%{
	if $input.Swigcptr() > 0 {
		$result = $input
	} else {
		$result = nil
	}
%}
%enddef
%typemap(goout) GDALDatasetShadow* CHECK_NULL()
%typemap(goout) GDALDriverShadow* CHECK_NULL()
