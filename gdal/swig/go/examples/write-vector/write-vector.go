package main

import (
	"flag"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/ogr"
	"io"
	"log"
	"os"
)

var drivername = flag.String("driver", "GeoJSON", "short driver name for output dataset")
var output = flag.String("output", "/vsistdout/out.geojson", "name for output dataset")
var layername = flag.String("layer", "point", "name for new point layer")
var fieldname = flag.String("field", "Name", "name for input field")

func main() {
	var driver gdal.Driver
	var ds gdal.Dataset
	var lyr ogr.Layer
	var fielddefn ogr.FieldDefn
	var err error

	flag.Parse()

	gdal.AllRegister()

	driver, err = gdal.GetDriverByName(*drivername)
	if err != nil {
		log.Fatal(err)
	}

	if ds, err = driver.Create(*output, 0, 0, 0, constant.GDT_Unknown); err != nil {
		log.Fatal(err)
	}
	defer ds.Close()

	if lyr, err = ds.CreateLayer(*layername, nil, ogr.WkbPoint, nil); err != nil {
		log.Fatal(err)
	}

	if fielddefn, err = ogr.NewFieldDefn(*fieldname, ogr.OFTString); err != nil {
		log.Fatal(err)
	}

	if err = fielddefn.SetWidth(32); err != nil {
		log.Fatal(err)
	}

	if err = lyr.CreateField(fielddefn, 1); err != nil {
		log.Fatal(err)
	}
	// TODO: destroy fielddefn

	for {
		var x, y float64
		var name string
		var n int

		n, err = fmt.Fscanf(os.Stdin, "%f,%f,%s\n", &x, &y, &name)
		if n == 0 && err == io.EOF {
			break // EOF
		} else if err != nil {
			log.Fatal(err)
		}

		var feat ogr.Feature
		if feat, err = ogr.NewFeature(lyr.GetLayerDefn()); err != nil {
			log.Fatal(err)
		}

		// TODO: use feat.GetFieldIndex(*fieldname) and SetFieldString
		if err = feat.SetField(*fieldname, name); err != nil {
			log.Fatal(err)
		}

		var pt ogr.Geometry
		if pt, err = ogr.NewGeometry(ogr.WkbPoint); err != nil {
			log.Fatal(err)
		}

		err = pt.SetPoint2D(0, x, y)
		if err != nil {
			log.Fatal(err)
		}

		if err = feat.SetGeometry(pt); err != nil {
			log.Fatal(err)
		}
		// TODO: destroy geometry

		if err = lyr.CreateFeature(feat); err != nil {
			log.Fatal(err)
		}
		// TODO: destroy feature
	}
}