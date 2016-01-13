package main

import (
	"errors"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/ogr"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errors.New("A vector dataset must be specified as a command line argument."))
	}
	input := os.Args[1]

	gdal.AllRegister()

	var ds gdal.Dataset
	var err error
	if ds, err = gdal.Open(input, constant.OF_VECTOR|constant.OF_READONLY|constant.OF_VERBOSE_ERROR); err != nil {
		log.Fatal(err)
	}

	driver := ds.GetDriver()
	fmt.Printf("Input driver is %s.\n", driver.GetLongName())

	lyr := ds.GetLayerByIndex(0)
	lyr.ResetReading()

	for feat := lyr.GetNextFeature(); feat != nil; feat = lyr.GetNextFeature() {
		featdefn := lyr.GetLayerDefn()

		for i := 0; i < featdefn.GetFieldCount(); i++ {
			fielddefn := featdefn.GetFieldDefn(i)

			fname := fielddefn.GetName()
			ftype := fielddefn.GetTypeName()

			switch fielddefn.GetType() {
			case ogr.OFTInteger:
				fallthrough
			case ogr.OFTInteger64:
				fmt.Printf("%s (%s): %d\n", fname, ftype, feat.GetFieldAsInteger(i))
			case ogr.OFTReal:
				fmt.Printf("%s (%s): %.3f\n", fname, ftype, feat.GetFieldAsDouble(i))
			case ogr.OFTString:
				fallthrough
			default:
				fmt.Printf("%s (%s): %s\n", fname, ftype, feat.GetFieldAsString(i))
			}
		}

		geom := feat.GetGeometryRef()
		if geom != nil && geom.GetGeometryType() == ogr.WkbPolygon {
			centroid := geom.Centroid()
			fmt.Printf("Centroid: %.3f, %.3f\n", centroid.GetX(), centroid.GetY())
		} else {
			fmt.Println("No polygon geometry")
		}
	}

	fmt.Println("Done")
}
