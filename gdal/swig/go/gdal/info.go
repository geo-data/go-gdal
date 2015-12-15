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

func NewInformer(options []string) (i Informer) {
	opts := newGDALInfoOptions(options)
	i = &info{
		opts,
	}

	runtime.SetFinalizer(i, func(i *info) {
		deleteGDALInfoOptions(i.options)
	})
	return
}

func (i *info) Info(ds Dataset) (result string, err error) {
	result = wrap_GDALInfo(ds, i.options)
	err = cpl.LastError()
	return
}
