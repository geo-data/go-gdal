package gdal

import (
	"errors"
	"fmt"
	"runtime"
)

type VectorTranslator interface {
	NameProcessor
}

type vtranslate struct {
	options  GDALVectorTranslateOptions
	datasets []Dataset
}

func NewVectorTranslator(options []string) (t VectorTranslator) {
	opts := newGDALVectorTranslateOptions(options)
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

	ds = wrapper_GDALVectorTranslateDestName(name, t.datasets[0], t.options)
	err = lastError()
	if ds != nil || err != nil {
		return
	}
	err = fmt.Errorf("Vector translate failed for %s", name)
	return
}

func (t *vtranslate) DestDS(ds Dataset) (err error) {
	if len(t.datasets) == 0 {
		err = errors.New("No input dataset")
		return
	}

	ok := wrapper_GDALVectorTranslateDestDS(ds, t.datasets[0], t.options)
	err = lastError()
	if ok == 1 || err != nil {
		return
	}
	err = errors.New("Vector translate failed for dataset")
	return
}
