package ogr

import (
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/ogr/errors"
)

func (l SwigcptrLayer) CreateFeature(f Feature) (err error) {
	var r int
	r, err = l.wrap_CreateFeature(f)
	if r != 0 && err == nil {
		err = errors.New(r, "Feature creation failed")
	}
	return
}

func (l SwigcptrLayer) CreateField(f FieldDefn, approxok int) (err error) {
	var r int
	r, err = l.wrap_CreateField(f, approxok)
	if r != 0 && err == nil {
		err = errors.New(r, "Field creation failed")
	}
	return
}
