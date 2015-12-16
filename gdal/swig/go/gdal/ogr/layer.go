package ogr

import (
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/ogr/errors"
)

func (l SwigcptrLayer) CreateFeature(f Feature) (err error) {
	defer cpl.ErrorTrap()(&err)
	if r := l.wrap_CreateFeature(f); r != 0 {
		err = errors.New(r, "Feature creation failed")
	}
	return
}

func (l SwigcptrLayer) CreateField(f FieldDefn, approxok int) (err error) {
	defer cpl.ErrorTrap()(&err)
	if r := l.wrap_CreateField(f, approxok); r != 0 {
		err = errors.New(r, "Field creation failed")
	}
	return
}
