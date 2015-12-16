package gdal

import (
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
)

func Open(filename string, args ...interface{}) (ret Dataset, err error) {
	defer cpl.ErrorTrap()(&err)
	ret = wrap_OpenEx(filename, args...)
	if ret == nil {
		err = fmt.Errorf("The dataset cannot be opened: %v", filename)
	}
	return
}

func GetDriverByName(name string) (driver Driver, err error) {
	defer cpl.ErrorTrap()(&err)
	driver = wrap_GetDriverByName(name)
	if driver == nil {
		err = fmt.Errorf("The driver is not registered: %v", name)
	}
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
	defer cpl.ErrorTrap()(&err)

	out := [6]float64{}
	gtout = out[:]
	if ok := wrap_InvGeoTransform(gtin[:], gtout); ok != 1 {
		err = fmt.Errorf("Non invertible transform: %v", gtin)
	}
	return
}
