package main

import (
	"flag"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
	"log"
	"os"
)

func usage(status int) {
	fmt.Fprintln(os.Stderr, `Usage: gdal-process <command> [<options>] [<args>]
where <command> is one of:
  translate
  vector-translate
  warp
  nearblack
  grid
  rasterize
  dem-hillshade
  dem-slope
  dem-aspect
  dem-color-relief
  dem-tri
  dem-tpi
  dem-roughness

e.g:
    gdal-process translate -options "-of VRT" input.tif output.vrt

Get help on available options by passing -help to the command e.g:
    gdal-process dem-color-relief -help

Description: gdal-process provides a common interface to gdal processing
commands which have been librarified according to RFC 59.1.  It is meant to
provide an example and test of their use in Go.
(https://trac.osgeo.org/gdal/wiki/rfc59.1_utilities_as_a_library).`)

	os.Exit(status)
}

type Processor func(options []string) (processor gdal.NameProcessor, err error)

type command struct {
	processor     Processor
	name          string
	dataset_count int
}

func (c *command) Processor(options []string) (processor gdal.NameProcessor, err error) {
	return c.processor(options)
}

func (c *command) Convert(options []string, input []gdal.Dataset, output string) (result gdal.Dataset, err error) {
	var processor gdal.NameProcessor
	processor, err = c.Processor(options)
	if err != nil {
		return
	}
	processor.SetDatasets(input)
	result, err = processor.DestName(output)
	return
}

func NewCommand(name string, processor Processor, dataset_count int) *command {
	return &command{
		processor,
		name,
		dataset_count,
	}
}

func main() {
	if len(os.Args) == 1 {
		usage(1)
	}

	var cmd *command
	subcommand := os.Args[1]
	action := flag.NewFlagSet(subcommand, flag.ExitOnError)
	var showInfo = true

	switch subcommand {
	case "translate":
		cmd = NewCommand("gdal_translate", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewTranslator(options)
		}, 2)
	case "vector-translate":
		cmd = NewCommand("ogr2ogr", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewVectorTranslator(options)
		}, 2)
		showInfo = false
	case "warp":
		cmd = NewCommand("gdalwarp", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewWarper(options)
		}, -1)
	case "nearblack":
		cmd = NewCommand("nearblack", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewNearblackProcessor(options)
		}, 2)
	case "grid":
		cmd = NewCommand("gdal_grid", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewGridProcessor(options)
		}, 2)
	case "rasterize":
		cmd = NewCommand("gdal_rasterize", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewRasterizer(options)
		}, 2)
	case "dem-hillshade":
		cmd = NewCommand("gdaldem hillshade", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewHillshadeProcessor(options)
		}, 2)
	case "dem-slope":
		cmd = NewCommand("gdaldem slope", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewSlopeProcessor(options)
		}, 2)
	case "dem-aspect":
		cmd = NewCommand("gdaldem aspect", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewAspectProcessor(options)
		}, 2)
	case "dem-color-relief":
		colorf := action.String("color-filename", "", "color text file")
		cmd = NewCommand("gdaldem color-relief", func(options []string) (np gdal.NameProcessor, err error) {
			if len(*colorf) == 0 {
				fmt.Fprintln(os.Stderr, "A color text file is required")
				usage(1)
			}
			var cr gdal.ColorReliefProcessor
			cr, err = gdal.NewColorReliefProcessor(options)
			if err != nil {
				return
			}
			cr.SetColorFilename(*colorf)
			np = cr
			return
		}, 2)
	case "dem-tri":
		cmd = NewCommand("gdaldem TRI", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewTRIProcessor(options)
		}, 2)
	case "dem-tpi":
		cmd = NewCommand("gdaldem TPI", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewTPIProcessor(options)
		}, 2)
	case "dem-roughness":
		cmd = NewCommand("gdaldem roughness", func(options []string) (gdal.NameProcessor, error) {
			return gdal.NewRoughnessProcessor(options)
		}, 2)
	default:
		fmt.Fprintf(os.Stderr, "%q is not valid command.\n", os.Args[1])
		usage(2)
	}

	var infoOpts string
	if showInfo {
		action.BoolVar(&showInfo, "info", true, "display information on output dataset")
		action.StringVar(&infoOpts, "info-options", "", "options for information output (as per gdalinfo)")
	}

	opts := action.String("options", "", fmt.Sprintf("%s options", cmd.name))
	action.Parse(os.Args[2:])

	switch cmd.dataset_count {
	case 2:
		if action.NArg() != 2 {
			fmt.Fprintln(os.Stderr, "An input and output dataset are required.")
			usage(1)
		}
	case -1:
		fallthrough
	default:
		if action.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "One or more input datasets and one output dataset are required.")
			usage(1)
		}

	}
	args := action.Args()

	// Get the GDAL specific options.
	options := cpl.ParseCommandLine(*opts)

	inputs := args[:len(args)-1]
	output := args[len(args)-1]

	datasets := make([]gdal.Dataset, len(inputs))

	gdal.AllRegister()

	for i, input := range inputs {
		ds, gerr := gdal.Open(input, constant.OF_READONLY|constant.OF_VERBOSE_ERROR)
		if gerr != nil {
			log.Fatal(gerr)
		}

		datasets[i] = ds
	}

	ds, err := cmd.Convert(options, datasets, output)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		ds.FlushCache()
		ds.Close()
	}()

	if showInfo {
		o := cpl.ParseCommandLine(infoOpts)
		i, err := gdal.NewInformer(o)
		if err != nil {
			log.Fatal(err)
		}

		info, err := i.Info(ds)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(info)
	}
}
