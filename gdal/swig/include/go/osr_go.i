
%insert("cgo_comment") %{
// "cpl_string.h" is needed to declare CSLDestroy.
#include "cpl_string.h"

// We need to link against GDAL when building the Go module.
#cgo LDFLAGS: -lgdal
%}
