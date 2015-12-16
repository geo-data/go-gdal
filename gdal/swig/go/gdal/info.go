package gdal

import (
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
	"runtime"
)

type Informer interface {
	Info(ds Dataset) (string, error)
}

type info struct {
	options GDALInfoOptions
}

func infoOptions(options []string) (opts GDALInfoOptions, err error) {
	cpl.ErrorReset()
	opts = newGDALInfoOptions(options)
	err = cpl.LastError()
	if err != nil && opts != nil {
		deleteGDALInfoOptions(opts)
	}
	return
}

func NewInformer(options []string) (i Informer, err error) {
	var opts GDALInfoOptions
	opts, err = infoOptions(options)
	if err != nil {
		return
	}

	i = &info{
		opts,
	}

	runtime.SetFinalizer(i, func(i *info) {
		deleteGDALInfoOptions(i.options)
	})
	return
}

func (i *info) Info(ds Dataset) (result string, err error) {
	defer cpl.ErrorTrap()(&err)
	result = wrap_GDALInfo(ds, i.options)
	return
}
