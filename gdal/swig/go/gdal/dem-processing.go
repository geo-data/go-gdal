package gdal

import (
	"errors"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
	"runtime"
)

type DEMProcessor interface {
	NameProcessor
	Processing() (processing string)
}

type ColorReliefProcessor interface {
	DEMProcessor
	ColorFilename() (filename string)
	SetColorFilename(filename string)
}

type dem struct {
	options               GDALDEMProcessingOptions
	datasets              []Dataset
	processing, colorfile string
}

func newDEM(options []string, processing string) (d *dem) {
	opts := newGDALDEMProcessingOptions(options)
	d = &dem{
		opts,
		[]Dataset{},
		processing,
		"",
	}

	runtime.SetFinalizer(d, func(d *dem) {
		deleteGDALDEMProcessingOptions(d.options)
	})
	return
}

func NewHillshadeProcessor(options []string) (h DEMProcessor) {
	return DEMProcessor(newDEM(options, "hillshade"))
}

func NewSlopeProcessor(options []string) (h DEMProcessor) {
	return DEMProcessor(newDEM(options, "slope"))
}

func NewAspectProcessor(options []string) (h DEMProcessor) {
	return DEMProcessor(newDEM(options, "aspect"))
}

func NewColorReliefProcessor(options []string) (h ColorReliefProcessor) {
	return ColorReliefProcessor(newDEM(options, "color-relief"))
}

func NewTRIProcessor(options []string) (h DEMProcessor) {
	return DEMProcessor(newDEM(options, "TRI"))
}

func NewTPIProcessor(options []string) (h DEMProcessor) {
	return DEMProcessor(newDEM(options, "TPI"))
}

func NewRoughnessProcessor(options []string) (h DEMProcessor) {
	return DEMProcessor(newDEM(options, "Roughness"))
}

func (d *dem) Processing() string {
	return d.processing
}

func (d *dem) ColorFilename() string {
	return d.colorfile
}

func (d *dem) SetColorFilename(colorfile string) {
	d.colorfile = colorfile
}

func (d *dem) Datasets() []Dataset {
	return d.datasets
}

func (d *dem) SetDatasets(datasets []Dataset) {
	d.datasets = datasets
}

func (d *dem) DestName(name string) (ds Dataset, err error) {
	if len(d.datasets) == 0 {
		err = errors.New("No input dataset")
		return
	}

	ds = wrapper_GDALDEMProcessing(name, d.datasets[0], d.processing, d.colorfile, d.options)
	err = cpl.LastError()
	if ds != nil || err != nil {
		return
	}
	err = fmt.Errorf("DEM %s processing failed for %s", d.processing, name)
	return
}
