// This is used by the tests to generate an in memory dataset.

package gdal_test

import (
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"testing"
)

// This needs to return the buffer as well to keep it in scope.
func memDataset(width uint, height uint) (buf []uint8, ds gdal.Dataset, err error) {
	buf = make([]uint8, width*height)
	connstr := fmt.Sprintf("MEM:::DATAPOINTER=%p,PIXELS=%d,LINES=%d,DATATYPE=Byte", &buf[0], width, height)
	ds, err = gdal.Open(connstr)
	return
}

func memDatasetTest(width uint, height uint, t *testing.T) (buf []uint8, ds gdal.Dataset) {
	var err error
	buf, ds, err = memDataset(width, height)
	if err != nil {
		t.Fatal(err)
	}
	return
}
