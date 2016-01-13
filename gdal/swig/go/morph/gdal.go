// This file is not designed to be compiled via go.  Instead it is read using
// the morph program which merges all interfaces found here with their
// corresponding interface in the appropriate package.
package gdal

type Dataset interface {
	GetGeoTransform() (gt []float64)
	SetGeoTransform(gt [6]float64) (err error)
}
