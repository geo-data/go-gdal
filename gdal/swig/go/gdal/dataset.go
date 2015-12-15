package gdal

import (
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/ogr"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/osr"
)

func (d SwigcptrDataset) GetGeoTransform() (gt []float64) {
	a := [6]float64{}
	gt = a[:]
	d.wrap_GetGeoTransform(gt)
	return
}

func (d SwigcptrDataset) SetGeoTransform(gt [6]float64) (err error) {
	if r := d.wrap_SetGeoTransform(gt[:]); r != 1 {
		err = cpl.LastError()
	}
	return
}

func (d SwigcptrDataset) GetLayerByName(name string) (lyr ogr.Layer, err error) {
	lyr = d.wrap_GetLayerByName(name)
	if lyr != nil {
		return
	}
	err = cpl.LastError()
	if err == nil {
		err = fmt.Errorf("The layer does not exist: %s", name)
	}
	return
}

func (d SwigcptrDataset) CreateLayer(name string, srs osr.SpatialReference, geomtype int, options []string) (lyr ogr.Layer, err error) {
	lyr = d.wrap_CreateLayer(name, srs, geomtype, options)
	if lyr != nil {
		return
	}
	err = cpl.LastError()
	if err == nil {
		err = fmt.Errorf("Layer creation failed: %s", name)
	}
	return
}
