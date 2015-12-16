package gdal

import "github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"

func (d SwigcptrDriver) CreateCopy(path string, ds Dataset, args ...interface{}) (ret Dataset, err error) {
	defer cpl.ErrorTrap()(&err)
	ret = d.wrap_CreateCopy(path, ds, args...)
	return
}

func (d SwigcptrDriver) Create(path string, xsize, ysize, bands, etype int, options []string) (ret Dataset, err error) {
	defer cpl.ErrorTrap()(&err)
	ret = d.wrap_Create(path, xsize, ysize, bands, etype, options)
	return
}
