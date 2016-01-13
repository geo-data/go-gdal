package gdal

import (
	"errors"
	"runtime"
)

const dserr = "No input datasets"

type Warper interface {
	DatasetProcessor
}

type warp struct {
	options  GDALWarpAppOptions
	datasets []Dataset
}

func warpAppOptions(options []string) (opts GDALWarpAppOptions, err error) {
	opts, err = newGDALWarpAppOptions(options)
	if err != nil && opts != nil {
		deleteGDALWarpAppOptions(opts)
	}
	return
}

func NewWarper(options []string) (w Warper, err error) {
	var opts GDALWarpAppOptions
	opts, err = warpAppOptions(options)
	if err != nil {
		return
	}

	w = &warp{
		opts,
		[]Dataset{},
	}

	runtime.SetFinalizer(w, func(w *warp) {
		deleteGDALWarpAppOptions(w.options)
	})
	return
}

func (w *warp) Datasets() []Dataset {
	return w.datasets
}

func (w *warp) SetDatasets(datasets []Dataset) {
	w.datasets = datasets
}

func (w *warp) DestName(name string) (ds Dataset, err error) {
	if len(w.datasets) == 0 {
		err = errors.New(dserr)
		return
	}

	return wrapper_GDALWarpDestName(name, len(w.datasets), w.datasets, w.options)
}

func (w *warp) DestDS(ds Dataset) (err error) {
	if len(w.datasets) == 0 {
		err = errors.New(dserr)
		return
	}

	var ok int
	ok, err = wrapper_GDALWarpDestDS(ds, len(w.datasets), w.datasets, w.options)
	if ok != 1 && err == nil {
		err = errors.New("Warp failed for dataset")
	}

	return
}
