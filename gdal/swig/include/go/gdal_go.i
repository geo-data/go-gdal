%import "typemaps_go.i"

%rename(wrap_CreateCopy) CreateCopy( const char *utf8_path, 
                                     GDALDatasetShadow* src, 
                                     int strict = 1, 
                                     char **options = 0, 
                                     GDALProgressFunc callback = NULL,
                                     void* callback_data = NULL);

%insert("cgo_comment") %{
// "gdal.h" is needed to declare GDALProgressFunc.
#include "gdal.h"

// "cpl_string.h" is needed to declare CSLDestroy.
#include "cpl_string.h"
  
// We need to link against GDAL when building the Go module.
#cgo LDFLAGS: -lgdal

// Declare the GDALProgressFunc we're exporting from Go.
extern int goGDALProgressFuncProxyA(double complete, char* message, void* data);
%}

