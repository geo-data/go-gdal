package gdal_test

import (
	"reflect"
	"testing"
)

func TestGetRasterXSize(t *testing.T) {
	_, ds := memDatasetTest(256, 200, t)
	x := ds.GetRasterXSize()
	if x != 256 {
		t.Errorf("Dataset.GetRasterXSize() == %v, expected 256", x)
	}
}

func TestGetRasterYSize(t *testing.T) {
	_, ds := memDatasetTest(200, 256, t)
	y := ds.GetRasterYSize()
	if y != 256 {
		t.Errorf("Dataset.GetRasterYSize() == %v, expected 256", y)
	}
}

func TestGetGeoTransform(t *testing.T) {
	_, ds := memDatasetTest(256, 256, t)
	expected := []float64{0, 1, 0, 0, 0, 1}
	gt := ds.GetGeoTransform()
	if !reflect.DeepEqual(gt, expected) {
		t.Errorf("Dataset.GetGeoTransform() == %v, does not match %v", gt, expected)
	}
}

func TestSetGeoTransform(t *testing.T) {
	_, ds := memDatasetTest(256, 256, t)
	good := [6]float64{600000.0, 9.0, 1.0, 4600000.0, 1.0, -20.0}
	err := ds.SetGeoTransform(good)
	if err != nil {
		t.Errorf("SetGeoTransform(%v): error is not nil: %s", good, err.Error())
	}
	gt := ds.GetGeoTransform()
	if !reflect.DeepEqual(gt, good[:]) {
		t.Errorf("Dataset.SetGeoTransform(%v) did not set the transform", good)
	}
}
