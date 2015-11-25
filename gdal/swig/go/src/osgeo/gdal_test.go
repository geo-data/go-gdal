package gdal

import (
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
	before = GetDriverCount()
	AllRegister()
	after = GetDriverCount()

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
	version := VersionInfo("--version")
	pattern := "^GDAL \\d+\\.\\d+\\.\\d+\\w*, released \\d{4}/\\d{2}/\\d{2}$"
	match, _ := regexp.MatchString(pattern, version)
	if !match {
		t.Errorf("VersionInfo(%v) == %v, does not match %v", "--version", version, pattern)
	}
}

func TestReadDir(t *testing.T) {
	dir := "."
	names, err := ReadDir(dir)
	if err != nil {
		t.Errorf("ReadDir(%v): error is not nil: %s", dir, err.Error())
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
	names, err = ReadDir(bad)
	if err == nil {
		t.Errorf("ReadDir(%v): bad directory but error is nil", bad)
	}
}

func TestOpen(t *testing.T) {
	bad := "/tmp/foobar"
	ds, err := Open(bad)
	if err == nil {
		t.Errorf("Open(%v): error is nil", bad)
	}
	if err.ErrorNo() != 4 {
		t.Errorf("Open(%v): error number == %v, expected 4", bad, err.ErrorNo())
	}
	if err.ErrorType() != 3 {
		t.Errorf("Open(%v): error type == %v, expected 3", bad, err.ErrorNo())
	}
	msg := err.Error()
	substr := "does not exist"
	if !strings.Contains(msg, substr) {
		t.Errorf("Open(%v): error message == %v, does not contain %v", bad, msg, substr)
	}
	if ds != nil {
		t.Errorf("Open(%v): datasource is not nil: %v", bad, ds)
	}

	good := "../../../../autotest/gcore/data/rgba.tif"
	ds, err = Open(good)
	if err != nil {
		t.Errorf("Open(%v): error is not nil: %v", good, err)
	}
	if ds == nil {
		t.Errorf("Open(%v): datasource is nil", good)
	}
}

func TestInvGeoTransform(t *testing.T) {
	bad := [6]float64{5, 3, 5, 3, 0, 0} // non invertible
	gtout, err := InvGeoTransform(bad)
	if err == nil {
		t.Errorf("InvGeoTransform(%v): non invertible but error is not nil", bad)
	}

	good := [6]float64{1, 2, 3, 4, 5, 6}
	gtout, err = InvGeoTransform(good)
	if err != nil {
		t.Errorf("InvGeoTransform(%v): error is not nil: %s", good, err.Error())
	}

	expected := []float64{-2, -2, 1, 1, 1.6666666666666665, -0.6666666666666666}
	if !reflect.DeepEqual(gtout, expected) {
		t.Errorf("InvGeoTransform(%v) == %v, does not match %v", good, gtout, expected)
	}
}
