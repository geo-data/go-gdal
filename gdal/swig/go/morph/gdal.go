// This file is not designed to be compiled via go.  Instead it is read using
// the morph program which merges all interfaces found here with their
// corresponding interface in the appropriate package.
package gdal

type Driver interface {
	Create(path string, xsize, ysize, bands, etype int, options []string) (ret Dataset, err error)
	CreateCopy(filename string, ds Dataset, args ...interface{}) (ret Dataset, err error)
}

type Dataset interface {
	CreateLayer(name string, srs osr.SpatialReference, geomtype int, options []string) (lyr ogr.Layer, err error)
	GetGeoTransform() (gt []float64)
	SetGeoTransform(gt [6]float64) (err error)
	GetLayerByName(name string) (lyr ogr.Layer, err error)
}
