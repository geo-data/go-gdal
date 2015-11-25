install: generate
	GOPATH=$$(pwd) go install -x osgeo

test: generate
	GOPATH=$$(pwd) go test osgeo

generate: ./src/osgeo/gdal.go

./src/osgeo/gdal.go: ./build/gdal.go ./src/osgeo/gdal_wrap.cpp ./merge/merge.go ./morph
	./morph

./build/gdal.go ./src/osgeo/gdal_wrap.cpp: ./build ../include/go/gdal_go.i ../include/gdal.i ../include/go/typemaps_go.i
	swig -Wall -I../include -I../include/go -outdir "./build" -I../../ -c++ -go -cgo -intgosize 64 -o src/osgeo/gdal_wrap.cpp ../include/gdal.i

./morph: src/morph/morph.go
	GOPATH=$$(pwd) go build src/morph/morph.go

./build:
	mkdir build

.PHONY: generate test install
