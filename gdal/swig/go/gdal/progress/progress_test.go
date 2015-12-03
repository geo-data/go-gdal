package gdal_test

import (
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/progress"
	"testing"
)

func TestProgress(t *testing.T) {
	_, ds := memDatasetTest(256, 200, t)

	mem_driver := gdal.GetDriverByName("MEM")
	if mem_driver == nil {
		t.Fatal("Cannot get MEM driver")
	}

	results := []string{}
	cb := progress.ProgressFunc(func(complete float64, message string, arg interface{}) int {
		results = append(results, arg.(string))
		return 1
	})

	arg := "testing normal progress."
	_, gerr := mem_driver.CreateCopy("/test.dat", ds, 0, []string{}, cb, arg)
	if gerr != nil {
		t.Fatalf("Cannot create copy MEM dataset: %s", gerr)
	}

	if len(results) < 3 {
		t.Errorf("CreateCopy() progress func len(results) == %d, expected >= 2", len(results))
	}

	for i, result := range results {
		if result != arg {
			t.Errorf("CreateCopy() progress func result %d == %s, expected %s", i, result, arg)
			break
		}
	}

	// Test for interrupt, interrupting after the progress function is called twice.
	results = []string{}
	cb = progress.ProgressFunc(func(complete float64, message string, arg interface{}) int {
		results = append(results, arg.(string))
		if len(results) < 2 {
			return 1
		}
		return 0
	})

	arg = "testing cancel progress."
	_, gerr = mem_driver.CreateCopy("/test2.dat", ds, 0, []string{}, cb, arg)
	if gerr == nil {
		t.Fatal("CreateCopy() err == nil on interrupt")
	} else {
		if gerr.ErrorType() != constant.CE_Failure {
			t.Errorf("CreateCopy() error type == %d, expected %d (CE_Failure)", gerr.ErrorType(), constant.CE_Failure)
		}
		if gerr.ErrorNo() != constant.CPLE_UserInterrupt {
			t.Errorf("CreateCopy() errorno == %d, expected %d (CPLE_UserInterrupt)", gerr.ErrorNo(), constant.CPLE_UserInterrupt)
		}
	}

	if len(results) != 2 {
		t.Errorf("CreateCopy() progress func len(results) == %d, expected 2", len(results))
	}

	for i, result := range results {
		if result != arg {
			t.Errorf("CreateCopy() progress func result %d == %s, expected %s", i, result, arg)
			break
		}
	}
}
