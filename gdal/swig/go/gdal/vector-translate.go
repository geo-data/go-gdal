package gdal

import (
	"errors"
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
	opts, err = newGDALVectorTranslateOptions(options)
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

	return wrapper_GDALVectorTranslateDestName(name, t.datasets[0], t.options)
}

func (t *vtranslate) DestDS(ds Dataset) (err error) {
	if len(t.datasets) == 0 {
		err = errors.New("No input dataset")
		return
	}

	var ok int
	ok, err = wrapper_GDALVectorTranslateDestDS(ds, t.datasets[0], t.options)
	if ok != 1 && err == nil {
		err = errors.New("Vector translate failed for dataset")
	}

	return
}
