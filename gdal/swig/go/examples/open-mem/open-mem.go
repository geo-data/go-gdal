package main

import (
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"log"
)

// memDataset creates an new buffer of type Byte along with a GDAL MEM dataset
// associated with the buffer.  The buffer is of size width * height, as is the
// raster.  Note that when the buffer goes out of scope and is garbage collected
// this invalidates the dataset which must no longer be used.
func memDataset(width uint, height uint) (buf []byte, ds gdal.Dataset, err error) {
	buf = make([]byte, width*height)
	connstr := fmt.Sprintf("MEM:::DATAPOINTER=%p,PIXELS=%d,LINES=%d,DATATYPE=Byte", &buf[0], width, height)
	ds, err = gdal.Open(connstr)
	return
}

func main() {
	gdal.AllRegister()

	width, height := uint(256), uint(128)
	log.Printf("Creating in memory raster of width %d x height %d", width, height)

	buf, ds, err := memDataset(width, height)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Buffer has length %d and capacity %d", len(buf), cap(buf))
	log.Printf("In memory raster has width %d and height %d", ds.GetRasterXSize(), ds.GetRasterYSize())
	log.Printf("Raster is of type %s", gdal.GetDataTypeName(ds.GetRasterBand(1).GetDataType()))
}
