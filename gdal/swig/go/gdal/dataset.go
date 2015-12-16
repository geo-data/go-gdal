package gdal

import (
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
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
	defer cpl.ErrorTrap()(&err)
	if r := d.wrap_SetGeoTransform(gt[:]); r != constant.CE_None {
		err = fmt.Errorf("SetGeoTransform(%v) failed: %d", gt, r)
	}
	return
}

func (d SwigcptrDataset) GetLayerByName(name string) (lyr ogr.Layer, err error) {
	defer cpl.ErrorTrap()(&err)
	lyr = d.wrap_GetLayerByName(name)
	return
}

func (d SwigcptrDataset) CreateLayer(name string, srs osr.SpatialReference, geomtype int, options []string) (lyr ogr.Layer, err error) {
	defer cpl.ErrorTrap()(&err)
	lyr = d.wrap_CreateLayer(name, srs, geomtype, options)
	if lyr == nil {
		err = fmt.Errorf("Layer creation failed: %s", name)
	}
	return
}
