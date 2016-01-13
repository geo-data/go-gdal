package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/progress"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var keep = flag.Bool("keep", false, "don't delete the output raster")

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal(errors.New("The source raster dataset must be specified as a command line argument."))
	}
	input := flag.Arg(0)

	gdal.AllRegister()

	var ds gdal.Dataset
	var err error
	if ds, err = gdal.Open(input, constant.OF_READONLY|constant.OF_VERBOSE_ERROR); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Opened input dataset (%dx%d).\n", ds.GetRasterXSize(), ds.GetRasterYSize())

	driver := ds.GetDriver()
	fmt.Printf("Input driver is %s.\n", driver.GetLongName())

	var dout string
	if dout, err = ioutil.TempDir("", "create-raster-copy-"); err != nil {
		log.Fatal(err)
	}

	// Delete the directory if required when the program exits.
	defer func() {
		if !*keep {
			log.Printf("Removing output directory %s", dout)
			if err := os.RemoveAll(dout); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// If something goes wrong, remove the temporary directory.
	kcache := *keep
	*keep = false

	var png_driver gdal.Driver
	if png_driver, err = gdal.GetDriverByName("PNG"); err != nil {
		log.Fatal(err)
	}

	count := 10000.0
	bar := pb.StartNew(int(count))
	bar.ShowCounters = false
	options := []string{"WORLDFILE=YES"}

	prog := progress.ProgressFunc(func(complete float64, message string, progressArg interface{}) int {
		bar.Set(int(count * complete))
		return 1
	})

	var png gdal.Dataset
	output := path.Join(dout, "copy.png")
	png, err = png_driver.CreateCopy(output, ds, 0, options, prog, "")
	bar.Finish()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("PNG dataset created (%dx%d): %s\n", png.GetRasterXSize(), png.GetRasterYSize(), output)

	// Reassign the the original choice of whether to remove the directory.
	*keep = kcache

	fmt.Println("Done")
}
