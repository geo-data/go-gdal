%import "typemaps_go.i"

%insert("cgo_comment") %{
// "cpl_string.h" is needed to declare CSLDestroy.
#include "cpl_string.h"

// We need to link against GDAL when building the Go module.
#cgo LDFLAGS: -lgdal
%}

%insert("header") %{
#include "cpl_vsi.h"
#include "cpl_multiproc.h"
#include "cpl_http.h"
%}

%include cpl.i
