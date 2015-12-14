package gdal

import (
	"errors"
	"fmt"
	"runtime"
)

type GridProcessor interface {
	NameProcessor
}

type grid struct {
	options  GDALGridOptions
	datasets []Dataset
}

func NewGridProcessor(options []string) (t GridProcessor) {
	opts := newGDALGridOptions(options)
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

	ds = wrapper_GDALGrid(name, t.datasets[0], t.options)
	err = lastError()
	if ds != nil || err != nil {
		return
	}
	err = fmt.Errorf("Grid failed for %s", name)
	return
}
