package ogr

import (
	"errors"
	ogrerr "github.com/geo-data/go-gdal/gdal/swig/go/gdal/ogr/errors"
)

func NewFeature(f FeatureDefn) (feat Feature, err error) {
	feat = wrap_NewFeature(f)
	//err = lastError()
	if err == nil && feat == nil {
		err = errors.New("Feature creation failed")
	}
	return
}

func (f SwigcptrFeature) SetField(a ...interface{}) (err error) {
	f.wrap_SetField(a...)
	return
}

func (f SwigcptrFeature) SetGeometry(geom Geometry) (err error) {
	r := f.wrap_SetGeometry(geom)
	// get error
	if err == nil && r != 0 {
		err = ogrerr.New(r, "Set geometry failed")
	}
	return
}
