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

func NewNearblackProcessor(options []string) (nb NearblackProcessor) {
	opts := newGDALNearblackOptions(options)
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
	ds = wrapper_GDALNearblackDestName(name, nb.datasets[0], nb.options)
	err = cpl.LastError()
	if ds != nil || err != nil {
		return
	}
	err = fmt.Errorf("Nearblack failed for %s", name)
	return
}

func (nb *nearblack) DestDS(ds Dataset) (err error) {
	if len(nb.datasets) == 0 {
		err = errors.New(dserr)
		return
	}
	ok := wrapper_GDALNearblackDestDS(ds, nb.datasets[0], nb.options)
	err = cpl.LastError()
	if ok == 1 || err != nil {
		return
	}
	err = errors.New("Nearblack failed for dataset")
	return
}
