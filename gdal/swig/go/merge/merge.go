// Package merge is not designed to be compiled via go.  Instead it is read
// using the morp program which merges all interfaces found here with their
// corresponding interface in the gdal package.
package merge

type Driver interface {
	CreateCopy(filename string, ds Dataset, args ...interface{}) (ret Dataset, err Err)
}
