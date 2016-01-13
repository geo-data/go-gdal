package gdal

import (
	"errors"
	"runtime"
)

type NearblackProcessor interface {
	DatasetProcessor
}

type nearblack struct {
	options  GDALNearblackOptions
	datasets []Dataset
}

func nearblackOptions(options []string) (opts GDALNearblackOptions, err error) {
	opts, err = newGDALNearblackOptions(options)
	if err != nil && opts != nil {
		deleteGDALNearblackOptions(opts)
	}
	return
}

func NewNearblackProcessor(options []string) (nb NearblackProcessor, err error) {
	var opts GDALNearblackOptions
	opts, err = nearblackOptions(options)
	if err != nil {
		return
	}

	nb = &nearblack{
		opts,
		[]Dataset{},
	}

	runtime.SetFinalizer(nb, func(nb *nearblack) {
		deleteGDALNearblackOptions(nb.options)
	})
	return
}

func (nb *nearblack) Datasets() []Dataset {
	return nb.datasets
}

func (nb *nearblack) SetDatasets(datasets []Dataset) {
	nb.datasets = datasets
}

func (nb *nearblack) DestName(name string) (ds Dataset, err error) {
	if len(nb.datasets) == 0 {
		err = errors.New(dserr)
		return
	}

	return wrapper_GDALNearblackDestName(name, nb.datasets[0], nb.options)
}

func (nb *nearblack) DestDS(ds Dataset) (err error) {
	if len(nb.datasets) == 0 {
		err = errors.New(dserr)
		return
	}

	var ok int
	ok, err = wrapper_GDALNearblackDestDS(ds, nb.datasets[0], nb.options)
	if ok != 1 && err == nil {
		err = errors.New("Nearblack failed for dataset")
	}

	return
}
