package gdal

import (
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
)

func (d SwigcptrDataset) GetGeoTransform() (gt []float64) {
	a := [6]float64{}
	gt = a[:]
	d.wrap_GetGeoTransform(gt)
	return
}

func (d SwigcptrDataset) SetGeoTransform(gt [6]float64) (err error) {
	var r int
	if r, err = d.wrap_SetGeoTransform(gt[:]); r != constant.CE_None && err == nil {
		err = fmt.Errorf("SetGeoTransform(%v) failed: %d", gt, r)
	}
	return
}
