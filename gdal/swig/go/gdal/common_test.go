// This is used by the tests to generate an in memory dataset.

package gdal_test

import (
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/raster/mem"
	"testing"
)

// This needs to return the buffer as well to keep it in scope.
func memDataset(width uint64, height uint64) (opts *mem.MemOptions, ds gdal.Dataset, err error) {
	opts, err = mem.New(width, height, constant.GDT_Byte)
	if err != nil {
		return
	}
	ds, err = opts.Dataset()
	return
}

func memDatasetTest(width uint64, height uint64, t *testing.T) (opts *mem.MemOptions, ds gdal.Dataset) {
	var err error
	opts, ds, err = memDataset(width, height)
	if err != nil {
		t.Fatal(err)
	}
	return
}
