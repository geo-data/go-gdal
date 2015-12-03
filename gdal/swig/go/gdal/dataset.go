package gdal

import "fmt"
import "github.com/geo-data/go-gdal/gdal/swig/go/gdal/ogr"

func (d SwigcptrDataset) GetGeoTransform() (gt []float64) {
	a := [6]float64{}
	gt = a[:]
	d.wrap_GetGeoTransform(gt)
	return
}

func (d SwigcptrDataset) SetGeoTransform(gt [6]float64) (err Err) {
	if r := d.wrap_SetGeoTransform(gt[:]); r != 1 {
		err = lastError()
	}
	return
}

func (d SwigcptrDataset) GetLayerByName(name string) (feat ogr.Layer, err error) {
	feat = d.wrap_GetLayerByName(name)
	if feat != nil {
		return
	}
	err = lastError()
	if err == nil {
		err = fmt.Errorf("The layer does not exist: %s", name)
	}
	return
}
