package ogr

import (
	"errors"
)

func NewGeometry(geomtype int) (geom Geometry, err error) {
	geom = wrap_NewGeometry(geomtype)
	//err = cpl.LastError()
	if err == nil && geom == nil {
		err = errors.New("Geometry creation failed")
	}
	return
}

func (g SwigcptrGeometry) SetPoint2D(i int, x, y float64) (err error) {
	g.wrap_SetPoint_2D(i, x, y)
	// Add error checking
	return
}
