package gdal_test

import (
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
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

	gerr, ok := err.(cpl.Error)
	if !ok {
		t.Fatalf("Open(%q) expected error type cpl.Error, got %T", bad, err)
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

func TestGetDataTypeName(t *testing.T) {
	var typetests = []struct {
		in  int
		out string
	}{
		{constant.GDT_Byte, "Byte"},
		{constant.GDT_UInt16, "UInt16"},
		{constant.GDT_Int16, "Int16"},
		{constant.GDT_UInt32, "UInt32"},
		{constant.GDT_Int32, "Int32"},
		{constant.GDT_Float32, "Float32"},
		{constant.GDT_Float64, "Float64"},
		{constant.GDT_CFloat32, "CFloat32"},
		{constant.GDT_CFloat64, "CFloat64"},
		{constant.GDT_Unknown, "Unknown"},
		{-1, ""},
	}

	for _, tt := range typetests {
		s := gdal.GetDataTypeName(tt.in)
		if s != tt.out {
			t.Errorf("GetDataTypeName(%q) => %q, want %q", tt.in, s, tt.out)
		}
	}
}
