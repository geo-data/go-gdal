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

func translateOptions(options []string) (opts GDALTranslateOptions, err error) {
	cpl.ErrorReset()
	opts = newGDALTranslateOptions(options)
	err = cpl.LastError()
	if err != nil && opts != nil {
		deleteGDALTranslateOptions(opts)
	}
	return
}

func NewTranslator(options []string) (t Translator, err error) {
	var opts GDALTranslateOptions
	opts, err = translateOptions(options)
	if err != nil {
		return
	}

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

	defer cpl.ErrorTrap()(&err)
	if ds = wrapper_GDALTranslate(name, t.datasets[0], t.options); ds == nil {
		err = fmt.Errorf("Translate failed for %s", name)
	}
	return
}
