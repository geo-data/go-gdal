package gdal

import (
	"testing"
)

var _ds Dataset

func getDs() Dataset {
	if _ds == nil {
		AllRegister()
		_ds, _ = Open("../../../../autotest/gcore/data/rgba.tif")
	}
	return _ds
}

func TestGetRasterXSize(t *testing.T) {
	ds := getDs()
	x := ds.GetRasterXSize()
	if x != 20 {
		t.Errorf("Dataset.GetRasterXSize() == %v, expected 20", x)
	}
}

func TestGetRasterYSize(t *testing.T) {
	ds := getDs()
	y := ds.GetRasterYSize()
	if y != 20 {
		t.Errorf("Dataset.GetRasterYSize() == %v, expected 20", y)
	}
}
