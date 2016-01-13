package gdal

import (
	"runtime"
)

type Informer interface {
	Info(ds Dataset) (string, error)
}

type info struct {
	options GDALInfoOptions
}

func infoOptions(options []string) (opts GDALInfoOptions, err error) {
	opts, err = newGDALInfoOptions(options)
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
	return wrap_GDALInfo(ds, i.options)
}
