package ogr

import (
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/ogr/errors"
)

func (l SwigcptrLayer) CreateFeature(f Feature) (err error) {
	r := l.wrap_CreateFeature(f)
	//err = cpl.LastError()
	if err == nil && r != 0 {
		err = errors.New(r, "Feature creation failed")
	}
	return
}

func (l SwigcptrLayer) CreateField(f FieldDefn, approxok int) (err error) {
	r := l.wrap_CreateField(f, approxok)
	//err = cpl.LastError()
	if err == nil && r != 0 {
		err = errors.New(r, "Field creation failed")
	}
	return
}
