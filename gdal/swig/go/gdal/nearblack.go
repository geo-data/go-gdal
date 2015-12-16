package gdal

import (
	"errors"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
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
	cpl.ErrorReset()
	opts = newGDALNearblackOptions(options)
	err = cpl.LastError()
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

	defer cpl.ErrorTrap()(&err)

	if ds = wrapper_GDALNearblackDestName(name, nb.datasets[0], nb.options); ds == nil {
		err = fmt.Errorf("Nearblack failed for %s", name)
	}

	return
}

func (nb *nearblack) DestDS(ds Dataset) (err error) {
	if len(nb.datasets) == 0 {
		err = errors.New(dserr)
		return
	}

	defer cpl.ErrorTrap()(&err)

	if ok := wrapper_GDALNearblackDestDS(ds, nb.datasets[0], nb.options); ok != 1 {
		err = errors.New("Nearblack failed for dataset")
	}

	return
}
