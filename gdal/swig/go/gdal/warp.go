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

func NewWarper(options []string) (w Warper) {
	opts := newGDALWarpAppOptions(options)
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
	ds = wrapper_GDALWarpDestName(name, len(w.datasets), w.datasets, w.options)
	err = cpl.LastError()
	if ds != nil || err != nil {
		return
	}
	err = fmt.Errorf("Warp failed for %s", name)
	return
}

func (w *warp) DestDS(ds Dataset) (err error) {
	if len(w.datasets) == 0 {
		err = errors.New(dserr)
		return
	}
	ok := wrapper_GDALWarpDestDS(ds, len(w.datasets), w.datasets, w.options)
	err = cpl.LastError()
	if ok == 1 || err != nil {
		return
	}
	err = errors.New("Warp failed for dataset")
	return
}
