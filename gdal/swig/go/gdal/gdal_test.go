package gdal_test

import (
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var before int
var after int

func TestMain(m *testing.M) {
	// Register drivers
	before = gdal.GetDriverCount()
	gdal.AllRegister()
	after = gdal.GetDriverCount()

	retCode := m.Run()
	os.Exit(retCode)
}

func TestAllRegister(t *testing.T) {
	if before != 0 {
		t.Errorf("GetDriverCount() == %v, does not match 0 before AllRegister()", before)
	}

	if after <= before {
		t.Errorf("GetDriverCount() == %v, <= %v after AllRegister()", after, before)
	}
}

func TestVersionInfo(t *testing.T) {
	version := gdal.VersionInfo("--version")
	pattern := "^GDAL \\d+\\.\\d+\\.\\d+\\w*, released \\d{4}/\\d{2}/\\d{2}$"
	match, _ := regexp.MatchString(pattern, version)
	if !match {
		t.Errorf("VersionInfo(%v) == %v, does not match %v", "--version", version, pattern)
	}
}

func TestReadDir(t *testing.T) {
	dir := "."
	names, err := gdal.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%v): error is not nil: %s", dir, err.Error())
	}
	if len(names) < 4 {
		t.Errorf("len(ReadDir(%v)) == %v, expected > 4", dir, len(names))
	}

	expected := []string{"gdal_test.go", "gdal.go", "gdal_wrap.cpp"}
Check:
	for _, name := range names {
		for i, e := range expected {
			if e != name {
				continue
			}
			expected = append(expected[:i], expected[i+1:]...) // delete the expected name
			if len(expected) == 0 {
				break Check
			}
		}
	}

	if len(expected) > 0 {
		t.Errorf("ReadDir(%v): expected files not found: %v", dir, expected)
	}

	bad := "oops"
	names, err = gdal.ReadDir(bad)
	if err == nil {
		t.Errorf("ReadDir(%v): bad directory but error is nil", bad)
	}
}

func TestOpenFail(t *testing.T) {
	bad := "nothing here"
	// Open without GDAL error
	ds, err := gdal.Open(bad, constant.OF_READONLY)
	if err == nil {
		t.Fatalf("Open(%q): error is nil", bad)
	}

	ds, err = gdal.Open(bad, constant.OF_READONLY|constant.OF_VERBOSE_ERROR)
	if err == nil {
		t.Fatalf("Open(%q): error is nil", bad)
	}

	gerr, ok := err.(gdal.Err)
	if !ok {
		t.Fatalf("Open(%q) expected error type Err, got %T", bad, err)
	}
	if gerr.ErrorNo() != 4 {
		t.Errorf("Open(%q): error number == %v, expected 4", bad, gerr.ErrorNo())
	}
	if gerr.ErrorType() != 3 {
		t.Errorf("Open(%q): error type == %v, expected 3", bad, gerr.ErrorNo())
	}
	msg := gerr.Error()
	substr := "does not exist"
	if !strings.Contains(msg, substr) {
		t.Errorf("Open(%q): error message == %v, does not contain %v", bad, msg, substr)
	}
	if ds != nil {
		t.Errorf("Open(%q): datasource is not nil: %v", bad, ds)
	}
}

func TestOpenSucceed(t *testing.T) {
	_, ds, err := memDataset(256, 256)
	if err != nil {
		t.Fatalf("Open(MEM): error is not nil: %v", err)
	}
	if ds == nil {
		t.Error("Open(MEM): datasource is nil")
	}
}

func TestInvGeoTransform(t *testing.T) {
	bad := [6]float64{5, 3, 5, 3, 0, 0} // non invertible
	gtout, err := gdal.InvGeoTransform(bad)
	if err == nil {
		t.Errorf("InvGeoTransform(%v): non invertible but error is not nil", bad)
	}

	good := [6]float64{1, 2, 3, 4, 5, 6}
	gtout, err = gdal.InvGeoTransform(good)
	if err != nil {
		t.Errorf("InvGeoTransform(%v): error is not nil: %s", good, err.Error())
	}

	expected := []float64{-2, -2, 1, 1, 1.6666666666666665, -0.6666666666666666}
	if !reflect.DeepEqual(gtout, expected) {
		t.Errorf("InvGeoTransform(%v) == %v, does not match %v", good, gtout, expected)
	}
}
