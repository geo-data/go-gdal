package ogr

import (
	"errors"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
)

func NewGeometry(geomtype int) (geom Geometry, err error) {
	defer cpl.ErrorTrap()(&err)
	if geom = wrap_NewGeometry(geomtype); geom == nil {
		err = errors.New("Geometry creation failed")
	}
	return
}

func (g SwigcptrGeometry) SetPoint2D(i int, x, y float64) (err error) {
	defer cpl.ErrorTrap()(&err)
	g.wrap_SetPoint_2D(i, x, y)
	return
}
