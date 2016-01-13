package ogr

import (
	ogrerr "github.com/geo-data/go-gdal/gdal/swig/go/gdal/ogr/errors"
)

func (f SwigcptrFeature) SetGeometry(geom Geometry) (err error) {
	var r int
	r, err = f.wrap_SetGeometry(geom)
	if err == nil && r != 0 {
		err = ogrerr.New(r, "Set geometry failed")
	}
	return
}
