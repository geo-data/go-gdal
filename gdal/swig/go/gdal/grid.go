package gdal

import (
	"errors"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
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
	cpl.ErrorReset()
	opts = newGDALGridOptions(options)
	err = cpl.LastError()
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

	defer cpl.ErrorTrap()(&err)

	if ds = wrapper_GDALGrid(name, t.datasets[0], t.options); ds == nil {
		err = fmt.Errorf("Grid failed for %s", name)
	}

	return
}
