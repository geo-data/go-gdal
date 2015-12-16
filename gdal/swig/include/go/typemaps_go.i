%include "typemaps.i"

%typemap(imtype) IF_FALSE_RETURN_NONE "int"

/* Replace NULL pointers in go strings with an empty string.  Declaring a string
   variable in go (var s string) produces a NULL pointer in C: we need to check
   for this and assign an empty string instead. */
%typemap(in) (const char * utf8_path) {
   $1 = (char *)$input.p;       /* From the default swig string typemap. */
   if (!$1) {
     $1 = (char *)"";
   }
}

%typemap(gotype) (const void *pBuffer) %{[]byte%}
%typemap(imtype) (const void *pBuffer) "*C.char"
%typemap(goin) (const void *pBuffer) %{
  $result = (*C.char)(unsafe.Pointer(&$input[0]))
%}
%apply (const void *pBuffer) {void *pBuffer};

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

%typemap(gotype) GDALDatasetShadow** poObjects %{[]Dataset%}
%typemap(imtype) GDALDatasetShadow** poObjects "unsafe.Pointer"
%typemap(goin) GDALDatasetShadow** poObjects %{
  	$input_l := len($input)
    $result_i := make([]unsafe.Pointer, $input_l)
	for i := 0; i < $input_l; i++ {
     $result_i[i] = unsafe.Pointer($input[i].Swigcptr())
	}
  $result = unsafe.Pointer(&$result_i[0])
%}

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
%typemap(gotype) ( GDALProgressFunc callback = NULL ) "progress.ProgressFunc"
%typemap(imtype) ( GDALProgressFunc callback = NULL ) "C.GDALProgressFunc"
%typemap(goin) ( GDALProgressFunc callback = NULL ) %{
    var progf progress.ProgressFunc
	if $input != nil {
        progf = $input
        $result = C.GDALProgressFunc(progress.GDALProgressFunc())
    }
%}
%typemap(goin) ( void* callback_data = NULL ) %{
    if progf != nil {
        $result = progress.New(progf, $input)
    }
%}
%typemap(in) ( GDALProgressFunc callback = NULL)
{
  $1 = $input;
}

%typemap(gotype) ( void* user_data ) "interface{}"
%typemap(imtype) ( void* user_data) "unsafe.Pointer"
%typemap(gotype) ( CPLErrorHandler pfnErrorHandler ) "ErrorHandler"
%typemap(imtype) ( CPLErrorHandler pfnErrorHandler ) "C.CPLErrorHandler"
%typemap(goin) ( CPLErrorHandler pfnErrorHandler ) %{
    var handler ErrorHandler
	if $input != nil {
        handler = $input
        $result = C.CPLErrorHandler(CPLErrorHandler())
    }
%}
%typemap(goin) ( void* user_data ) %{
    if handler != nil {
        $result = NewErrorHandler(handler, $input)
    }
%}
%typemap(in) ( CPLErrorHandler pfnErrorHandler)
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
%typemap(goout) OGRLayerShadow* CHECK_NULL()
%typemap(goout) OGRFeatureShadow* CHECK_NULL()

/* Allow srs objects to be NULL. */
%typemap(imtype) OSRSpatialReferenceShadow* srs "uintptr"
%typemap(goin) OSRSpatialReferenceShadow* srs %{
  if $input != nil {
      $result = $input.Swigcptr()
  }
%}
