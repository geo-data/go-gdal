// http://www.gdal.org/frmt_mem.html

package mem

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"unsafe"
)

type MemOptions struct {
	Buffer                              []byte
	Pixels, Lines                       uint64
	Bands                               uint
	Type                                int // Gdal data type
	PixelOffset, LineOffset, BandOffset uint
}

// MemOptions needs to be in scope for as long as any returned gdal.Dataset.
// This is because MemOptions.Buffer is shared by the GDAL memory dataset and
// must not be garbage collected before the dataset is closed.
func (m *MemOptions) Dataset() (dataset gdal.Dataset, err error) {
	if len(m.Buffer) == 0 {
		err = errors.New("Memory buffer is empty")
		return
	}
	dtn := gdal.GetDataTypeName(m.Type)

	connstr := fmt.Sprintf("MEM:::DATAPOINTER=%p,PIXELS=%d,LINES=%d,DATATYPE=%s,BANDS=%d,PIXELOFFSET=%d,LINE0FFSET=%d,BANDOFFSET=%d", &m.Buffer[0], m.Pixels, m.Lines, dtn, m.Bands, m.PixelOffset, m.LineOffset, m.BandOffset)
	dataset, err = gdal.Open(connstr, constant.OF_VERBOSE_ERROR)
	return
}

func New(width, height uint64, datatype int) (options *MemOptions, err error) {
	buf := new(bytes.Buffer)

	// Get the native byte order
	// (https://groups.google.com/d/msg/golang-nuts/zmh64YkqOV8/iJe-TrTTeREJ).
	var order binary.ByteOrder
	var x uint32 = 0x01020304
	switch *(*byte)(unsafe.Pointer(&x)) {
	case 0x01:
		order = binary.BigEndian
	case 0x04:
		order = binary.LittleEndian
	}

	switch datatype {
	case constant.GDT_Byte:
		binary.Write(buf, order, make([]byte, width*height))
	case constant.GDT_UInt16:
		binary.Write(buf, order, make([]uint16, width*height))
	case constant.GDT_Int16:
		binary.Write(buf, order, make([]int16, width*height))
	case constant.GDT_UInt32:
		binary.Write(buf, order, make([]uint32, width*height))
	case constant.GDT_Int32:
		binary.Write(buf, order, make([]int32, width*height))
	case constant.GDT_Float32:
		binary.Write(buf, order, make([]float32, width*height))
	case constant.GDT_Float64:
		binary.Write(buf, order, make([]float64, width*height))
	case constant.GDT_CFloat32:
		binary.Write(buf, order, make([]complex64, width*height))
	case constant.GDT_CFloat64:
		binary.Write(buf, order, make([]complex128, width*height))
	case constant.GDT_Unknown:
		fallthrough
	case constant.GDT_CInt16:
		fallthrough
	case constant.GDT_CInt32:
		fallthrough
	default:
		name := gdal.GetDataTypeName(datatype)
		if len(name) > 0 {
			err = fmt.Errorf("Unsupported GDAL datatype: %s", gdal.GetDataTypeName(datatype))
		} else {
			err = fmt.Errorf("Unhandled GDAL datatype: %d", datatype)
		}
		return
	}

	options = &MemOptions{
		buf.Bytes(),
		width, height,
		1,
		datatype,
		0, 0, 0,
	}
	return
}
