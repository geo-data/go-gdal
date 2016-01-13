package gdal

import (
	"fmt"
)

func Open(filename string, args ...interface{}) (ret Dataset, err error) {
	ret, err = wrap_OpenEx(filename, args...)
	return
}

func ApplyGeoTransform(transform [6]float64, pixel float64, line float64) (geox float64, geoy float64) {
	x := []float64{0.0}
	y := []float64{0.0}
	wrap_ApplyGeoTransform(transform[:], pixel, line, x, y)
	geox, geoy = x[0], y[0]
	return
}

func InvGeoTransform(gtin [6]float64) (gtout []float64, err error) {
	out := [6]float64{}
	gtout = out[:]
	var ok int
	ok, err = wrap_InvGeoTransform(gtin[:], gtout)
	if err == nil && ok != 1 {
		err = fmt.Errorf("Non invertible transform: %v", gtin)
	}
	return
}
