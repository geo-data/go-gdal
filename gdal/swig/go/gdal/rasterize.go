package gdal

import (
	"errors"
	"runtime"
)

type Rasterizer interface {
	DatasetProcessor
}

type rasterize struct {
	options  GDALRasterizeOptions
	datasets []Dataset
}

func rasterizeOptions(options []string) (opts GDALRasterizeOptions, err error) {
	opts, err = newGDALRasterizeOptions(options)
	if err != nil && opts != nil {
		deleteGDALRasterizeOptions(opts)
	}
	return
}

func NewRasterizer(options []string) (r Rasterizer, err error) {
	var opts GDALRasterizeOptions
	opts, err = rasterizeOptions(options)
	if err != nil {
		return
	}

	r = &rasterize{
		opts,
		[]Dataset{},
	}

	runtime.SetFinalizer(r, func(r *rasterize) {
		deleteGDALRasterizeOptions(r.options)
	})
	return
}

func (r *rasterize) Datasets() []Dataset {
	return r.datasets
}

func (r *rasterize) SetDatasets(datasets []Dataset) {
	r.datasets = datasets
}

func (r *rasterize) DestName(name string) (ds Dataset, err error) {
	if len(r.datasets) == 0 {
		err = errors.New(dserr)
		return
	}

	return wrapper_GDALRasterizeDestName(name, r.datasets[0], r.options)
}

func (r *rasterize) DestDS(ds Dataset) (err error) {
	if len(r.datasets) == 0 {
		err = errors.New(dserr)
		return
	}

	var ok int
	ok, err = wrapper_GDALRasterizeDestDS(ds, r.datasets[0], r.options)
	if ok != 1 && err == nil {
		err = errors.New("Rasterize failed for dataset")
	}

	return
}
