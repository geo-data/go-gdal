package cpl

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
	var h uintptr
	h, err = VSIFOpenL(name, access)
	if h == 0 && err == nil {
		err = fmt.Errorf("Cannot open file %v using %v", name, access)
		return
	}
	visf = vsi(h)
	return
}

func (v vsi) Tell() (o int64, err error) {
	return VSIFTellL(uintptr(v))
}

func (v vsi) Seek(offset int64, whence int) (noffset int64, err error) {
	var ret int
	ret, err = VSIFSeekL(uintptr(v), offset, whence)
	if ret != 0 && err == nil {
		err = errors.New("Seek failed")
		return
	}
	noffset, err = v.Tell()
	return
}

func (v vsi) Read(p []byte) (n int, err error) {
	l := len(p)
	up := uintptr(v)

	var nl int64
	nl, err = VSIFReadL(p, 1, int64(l), up)
	if err != nil {
		return
	}
	n = int(nl)

	var eol int
	eol, err = VSIFEofL(up)
	if err != nil {
		return
	}

	if eol == 1 {
		err = io.EOF
	}
	return
}

func (v vsi) Write(p []byte) (n int, err error) {
	l := len(p)
	var nl int64
	nl, err = VSIFWriteL(p, 1, int64(l), uintptr(v))
	if err != nil {
		return
	}
	n = int(nl)

	if n != l {
		err = fmt.Errorf("Write failed, only wrote %d of %d bytes", n, l)
	}
	return
}

func (v vsi) Close() (err error) {
	var r int
	r, err = VSIFCloseL(uintptr(v))
	if err != nil {
		return
	}

	if r != 0 {
		return errors.New("Close failed")
	}

	return
}

func (v vsi) Truncate(size int64) (err error) {
	var r int
	r, err = VSIFTruncateL(uintptr(v), size)
	if err != nil {
		return
	}

	if r != 0 {
		err = errors.New("Truncate failed")
	}

	return
}
