package gdal

import (
	"errors"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
	"runtime"
)

type VectorTranslator interface {
	NameProcessor
}

type vtranslate struct {
	options  GDALVectorTranslateOptions
	datasets []Dataset
}

func vectorTranslateOptions(options []string) (opts GDALVectorTranslateOptions, err error) {
	cpl.ErrorReset()
	opts = newGDALVectorTranslateOptions(options)
	err = cpl.LastError()
	if err != nil && opts != nil {
		deleteGDALVectorTranslateOptions(opts)
	}
	return
}

func NewVectorTranslator(options []string) (t VectorTranslator, err error) {
	var opts GDALVectorTranslateOptions
	opts, err = vectorTranslateOptions(options)
	if err != nil {
		return
	}

	t = &vtranslate{
		opts,
		[]Dataset{},
	}

	runtime.SetFinalizer(t, func(t *vtranslate) {
		deleteGDALVectorTranslateOptions(t.options)
	})
	return
}

func (t *vtranslate) Datasets() []Dataset {
	return t.datasets
}

func (t *vtranslate) SetDatasets(datasets []Dataset) {
	t.datasets = datasets
}

func (t *vtranslate) DestName(name string) (ds Dataset, err error) {
	if len(t.datasets) == 0 {
		err = errors.New("No input dataset")
		return
	}

	defer cpl.ErrorTrap()(&err)

	if ds = wrapper_GDALVectorTranslateDestName(name, t.datasets[0], t.options); ds == nil {
		err = fmt.Errorf("Vector translate failed for %s", name)
	}

	return
}

func (t *vtranslate) DestDS(ds Dataset) (err error) {
	if len(t.datasets) == 0 {
		err = errors.New("No input dataset")
		return
	}

	defer cpl.ErrorTrap()(&err)

	if ok := wrapper_GDALVectorTranslateDestDS(ds, t.datasets[0], t.options); ok != 1 {
		err = errors.New("Vector translate failed for dataset")
	}

	return
}
