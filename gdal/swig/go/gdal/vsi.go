package gdal

import (
	"errors"
	"fmt"
	"io"
)

type vsi uintptr

type VSIF interface {
	io.ReadWriteCloser
	io.Seeker
	Tell() (int64, error)
	Truncate(size int64) (err error)
}

func VSIFOpen(name, access string) (visf VSIF, err error) {
	r := vsi(VSIFOpenL(name, access))
	err = lastError()
	if err != nil {
		return
	}
	if r == 0 {
		err = fmt.Errorf("Cannot open file %v using %v", name, access)
		return
	}
	visf = r
	return
}

func (v vsi) Tell() (o int64, err error) {
	o = VSIFTellL(uintptr(v))
	err = lastError()
	return
}

func (v vsi) Seek(offset int64, whence int) (noffset int64, err error) {
	ret := VSIFSeekL(uintptr(v), offset, whence)
	err = lastError()
	if err != nil {
		return
	}
	if ret != 0 {
		err = errors.New("Seek failed")
		return
	}
	noffset, err = v.Tell()
	return
}

func (v vsi) Read(p []byte) (n int, err error) {
	l := len(p)
	up := uintptr(v)
	n = int(VSIFReadL(p, 1, int64(l), up))
	err = lastError()
	if err != nil {
		return
	}
	if VSIFEofL(up) == 1 {
		err = io.EOF
	}
	return
}

func (v vsi) Write(p []byte) (n int, err error) {
	l := len(p)
	n = int(VSIFWriteL(p, 1, int64(l), uintptr(v)))
	err = lastError()
	if err != nil {
		return
	}
	if n != l {
		err = fmt.Errorf("Write failed, only wrote %d of %d bytes", n, l)
	}
	return
}

func (v vsi) Close() error {
	r := VSIFCloseL(uintptr(v))
	err := lastError()
	if err != nil {
		return err
	}
	if r != 0 {
		return errors.New("Close failed")
	}
	return nil
}

func (v vsi) Truncate(size int64) (err error) {
	r := VSIFTruncateL(uintptr(v), size)
	err = lastError()
	if err != nil {
		return
	}
	if r != 0 {
		err = errors.New("Truncate failed")
	}
	return
}
