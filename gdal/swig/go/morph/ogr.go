// This file is not designed to be compiled via go.  Instead it is read using
// the morph program which merges all interfaces found here with their
// corresponding interface in the appropriate package.
package ogr

type Layer interface {
	CreateField(f FieldDefn, approxok int) (err error)
	CreateFeature(f Feature) (err error)
}

type Feature interface {
	SetField(a ...interface{}) error
	SetGeometry(geom Geometry) (err error)
}

type Geometry interface {
	SetPoint2D(i int, x, y float64) (err error)
}
