package gdal

import (
	"errors"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
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
	cpl.ErrorReset()
	opts = newGDALWarpAppOptions(options)
	err = cpl.LastError()
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

	defer cpl.ErrorTrap()(&err)

	if ds = wrapper_GDALWarpDestName(name, len(w.datasets), w.datasets, w.options); ds == nil {
		err = fmt.Errorf("Warp failed for %s", name)
	}

	return
}

func (w *warp) DestDS(ds Dataset) (err error) {
	if len(w.datasets) == 0 {
		err = errors.New(dserr)
		return
	}

	defer cpl.ErrorTrap()(&err)

	if ok := wrapper_GDALWarpDestDS(ds, len(w.datasets), w.datasets, w.options); ok == 1 {
		err = errors.New("Warp failed for dataset")
	}

	return
}
