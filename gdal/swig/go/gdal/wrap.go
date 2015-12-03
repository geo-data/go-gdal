package gdal

import "fmt"

func Open(filename string, args ...interface{}) (ret Dataset, err error) {
	ret = wrap_OpenEx(filename, args...)
	if ret != nil {
		return
	}
	err = lastError()
	if err == nil {
		err = fmt.Errorf("The dataset cannot be opened: %v", filename)
	}
	return
}

func ReadDir(utf8_path string) (ret []string, err error) {
	ret = wrap_ReadDir(utf8_path)
	if len(ret) > 0 {
		return
	}
	err = fmt.Errorf("Cannot read directory %v", utf8_path)
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
	ok := wrap_InvGeoTransform(gtin[:], gtout)
	if ok != 1 {
		err = fmt.Errorf("Non invertible transform: %v", gtin)
	}
	return
}
