package gdal_test

import (
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/progress"
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

func TestGetFileList(t *testing.T) {
	geojson := `{
"type": "FeatureCollection",
"features": [
{ "type": "Feature", "properties": { "FID": 1, "NAME": "Point 1" }, "geometry": { "type": "Point", "coordinates": [ 100.0, 0.0 ] } }
]
}`
	output := "/vsimem/out.shp"

	src, err := gdal.Open(geojson, constant.OF_READONLY|constant.OF_VERBOSE_ERROR)
	if err != nil {
		t.Fatalf("Open(geojson): error is not nil: %v", err)
	}
	defer src.Close()

	files := src.GetFileList()
	if len(files) != 0 {
		t.Errorf("len(GetFileList()) == %d, expected 0", len(files))
	}

	var shp gdal.Driver
	shp, err = gdal.GetDriverByName("ESRI Shapefile")
	if err != nil {
		t.Fatalf("GetDriverByName(\"ESRI Shapefile\"): error is not nil: %v", err)
	}

	prog := progress.ProgressFunc(func(complete float64, message string, progressArg interface{}) int {
		return 1
	})

	var dst gdal.Dataset
	dst, err = shp.CreateCopy(output, src, 0, []string{}, prog, "")
	if err != nil {
		t.Fatalf("CreateCopy(%q): error is not nil: %v", output, err)
	}
	defer dst.Close()

	files = dst.GetFileList()
	expected := []string{"/vsimem/out.shp", "/vsimem/out.shx", "/vsimem/out.dbf", ""}

	if !reflect.DeepEqual(files, expected) {
		t.Errorf("GetFileList() == %v, does not match %v", files, expected)
	}
}
