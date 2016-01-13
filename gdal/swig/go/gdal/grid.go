package gdal

import (
	"errors"
	"runtime"
)

type GridProcessor interface {
	NameProcessor
}

type grid struct {
	options  GDALGridOptions
	datasets []Dataset
}

func gridOptions(options []string) (opts GDALGridOptions, err error) {
	opts, err = newGDALGridOptions(options)
	if err != nil && opts != nil {
		deleteGDALGridOptions(opts)
	}
	return
}

func NewGridProcessor(options []string) (t GridProcessor, err error) {
	var opts GDALGridOptions
	opts, err = gridOptions(options)
	if err != nil {
		return
	}

	t = &grid{
		opts,
		[]Dataset{},
	}

	runtime.SetFinalizer(t, func(t *grid) {
		deleteGDALGridOptions(t.options)
	})
	return
}

func (t *grid) Datasets() []Dataset {
	return t.datasets
}

func (t *grid) SetDatasets(datasets []Dataset) {
	t.datasets = datasets
}

func (t *grid) DestName(name string) (ds Dataset, err error) {
	if len(t.datasets) == 0 {
		err = errors.New("No input dataset")
		return
	}

	return wrapper_GDALGrid(name, t.datasets[0], t.options)
}
