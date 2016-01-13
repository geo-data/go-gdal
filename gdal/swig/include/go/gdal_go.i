%go_import("github.com/geo-data/go-gdal/gdal/swig/go/gdal/progress")
%go_import("github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl")

%import "typemaps_go.i"

%insert("cgo_comment") %{
// "gdal.h" is needed to declare GDALProgressFunc.
#include "gdal.h"

// "cpl_string.h" is needed to declare CSLDestroy.
#include "cpl_string.h"

// We need to link against GDAL when building the Go module.
#cgo LDFLAGS: -lgdal
%}

#define FROM_GDAL_I
%import ogr.i
