package main

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"log"
	"osgeo"
)

func progress(complete float64, message string, progressArg interface{}) int {
	fmt.Printf("woo hoo here: %v %s\n", complete, progressArg.(string))
	return 1
}

func main() {
	fmt.Printf("Hello, world.\n")
	info := gdal.VersionInfo("--version")
	fmt.Printf("Version is %s.\n", info)
	fmt.Printf("Driver count before registering: %d\n", gdal.GetDriverCount())
	gdal.AllRegister()
	fmt.Printf("Driver count after registering: %d\n", gdal.GetDriverCount())

	/*gdal.Debug("trial", "setting CPLQuietErrorHandler")
	if err := gdal.SetErrorHandler("CPLQuietErrorHandler"); err != nil {
		log.Fatal(err)
	}*/

	geox, geoy := gdal.ApplyGeoTransform([6]float64{444720, 30, 0, 3751320, 0, -30}, 5.0, 10.0)
	fmt.Printf("ApplyGeoTransform() x = %v, y = %v\n", geox, geoy)

	if gtout, err := gdal.InvGeoTransform([6]float64{5, 3, 5, 3, 0, 0}); err == nil { // non invertible
		fmt.Printf("InvGeoTransform() t = %v\n", gtout)
	} else {
		log.Fatal(err)
	}

	dir, err := gdal.ReadDir("/tmp")
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range dir {
		fmt.Printf("Dir: %s\n", val)
	}

	//ds := gdal.Open("/tmp/gdal/autotest/gcore/data/rgba.tif")
	ds, gerr := gdal.Open("/tmp/gdal/originimmob.tif")
	if gerr == nil {
		fmt.Printf("Opened input file (%dx%d).\n", ds.GetRasterXSize(), ds.GetRasterYSize())
		driver := ds.GetDriver()
		fmt.Printf("Driver is %s.\n", driver.GetLongName())
		fmt.Printf("Dataset projection: %q\n", ds.GetProjection())

		png_driver := gdal.GetDriverByName("PNG")
		//png_driver := gdal.GetDriverByName("AAIGrid")
		if png_driver != nil {
			count := 10000.0
			bar := pb.StartNew(int(count))
			bar.ShowCounters = false
			png_name := "/tmp/rgba.png"
			options := []string{"WORLDFILE=YES"}

			prog := func(complete float64, message string, progressArg interface{}) int {
				bar.Set(int(count * complete))
				return 1
			}

			png, gerr := png_driver.CreateCopy(png_name, ds, 0, options, gdal.ProgressFunc(prog), "")
			bar.Finish()
			if gerr != nil {
				fmt.Printf("Error creating copy: %d, %s, %d\n", gerr.ErrorNo(), gerr.Error(), gerr.ErrorType())
			} else {
				fmt.Printf("PNG dataset created (%dx%d): %s\n", png.GetRasterXSize(), png.GetRasterYSize(), png_name)
			}
		} else {
			fmt.Printf("The driver could not be found.\n")
		}
	} else {
		fmt.Printf("Error: %d, %s, %d\n", gerr.ErrorNo(), gerr.Error(), gerr.ErrorType())
	}
	fmt.Printf("Goodbye, world.\n")
}
