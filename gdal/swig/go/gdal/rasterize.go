package gdal

import (
	"errors"
	"fmt"
	"runtime"
)

type Rasterizer interface {
	DatasetProcessor
}

type rasterize struct {
	options  GDALRasterizeOptions
	datasets []Dataset
}

func NewRasterizer(options []string) (r Rasterizer) {
	opts := newGDALRasterizeOptions(options)
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
	ds = wrapper_GDALRasterizeDestName(name, r.datasets[0], r.options)
	err = lastError()
	if ds != nil || err != nil {
		return
	}
	err = fmt.Errorf("Rasterize failed for %s", name)
	return
}

func (r *rasterize) DestDS(ds Dataset) (err error) {
	if len(r.datasets) == 0 {
		err = errors.New(dserr)
		return
	}
	ok := wrapper_GDALRasterizeDestDS(ds, r.datasets[0], r.options)
	err = lastError()
	if ok == 1 || err != nil {
		return
	}
	err = errors.New("Rasterize failed for dataset")
	return
}
