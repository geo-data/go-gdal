package gdal

import "github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"

func (d SwigcptrDriver) CreateCopy(path string, ds Dataset, args ...interface{}) (ret Dataset, err error) {
	ret = d.wrap_CreateCopy(path, ds, args...)
	if ret != nil {
		return
	}
	err = cpl.LastError()
	ret = nil
	return
}

func (d SwigcptrDriver) Create(path string, xsize, ysize, bands, etype int, options []string) (ret Dataset, err error) {
	ret = d.wrap_Create(path, xsize, ysize, bands, etype, options)
	if ret != nil {
		return
	}
	err = cpl.LastError()
	ret = nil
	return
}
