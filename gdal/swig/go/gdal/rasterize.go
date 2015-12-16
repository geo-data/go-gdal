package gdal

import (
	"errors"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
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
	cpl.ErrorReset()
	opts = newGDALRasterizeOptions(options)
	err = cpl.LastError()
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

	defer cpl.ErrorTrap()(&err)

	if ds = wrapper_GDALRasterizeDestName(name, r.datasets[0], r.options); ds == nil {
		err = fmt.Errorf("Rasterize failed for %s", name)
	}

	return
}

func (r *rasterize) DestDS(ds Dataset) (err error) {
	if len(r.datasets) == 0 {
		err = errors.New(dserr)
		return
	}

	defer cpl.ErrorTrap()(&err)

	if ok := wrapper_GDALRasterizeDestDS(ds, r.datasets[0], r.options); ok != 1 {
		err = errors.New("Rasterize failed for dataset")
	}

	return
}
