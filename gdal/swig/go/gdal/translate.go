package gdal

import (
	"errors"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
	"runtime"
)

type Translator interface {
	NameProcessor
}

type translate struct {
	options  GDALTranslateOptions
	datasets []Dataset
}

func NewTranslator(options []string) (t Translator) {
	opts := newGDALTranslateOptions(options)
	t = &translate{
		opts,
		[]Dataset{},
	}

	runtime.SetFinalizer(t, func(t *translate) {
		deleteGDALTranslateOptions(t.options)
	})
	return
}

func (t *translate) Datasets() []Dataset {
	return t.datasets
}

func (t *translate) SetDatasets(datasets []Dataset) {
	t.datasets = datasets
}

func (t *translate) DestName(name string) (ds Dataset, err error) {
	if len(t.datasets) == 0 {
		err = errors.New("No input dataset")
		return
	}

	ds = wrapper_GDALTranslate(name, t.datasets[0], t.options)
	err = cpl.LastError()
	if ds != nil || err != nil {
		return
	}
	err = fmt.Errorf("Translate failed for %s", name)
	return
}
