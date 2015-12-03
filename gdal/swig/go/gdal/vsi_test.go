package gdal

import (
	"encoding/csv"
	"io"
	"reflect"
	"testing"
)

func TestVSIF(t *testing.T) {
	// Test Open
	fname := "/vsimem/test.dat"
	fp, err := VSIFOpen(fname, "rw")
	if err != nil {
		t.Errorf("VSIFOpen(%v) error is not nil: %s", fname, err)
	}

	// Test Write
	in := []byte("this is a test.")
	var c int
	c, err = fp.Write(in)
	if err != nil {
		t.Errorf("VSIF.Write(%v) error is not nil: %s", in, err)
	}
	if c != len(in) {
		t.Errorf("VSIF.Write(%v) failed, only wrote %d of %d bytes", in, c, len(in))
	}

	// Test Tell
	var offset int64
	offset, err = fp.Tell()
	if err != nil {
		t.Errorf("VSIF.Tell() error is not nil: %s", err)
	}
	if offset != int64(c) {
		t.Errorf("VSIF.Tell() == %d, expected %d", offset, c)
	}

	// Test Seek
	offset, err = fp.Seek(0, 0)
	if err != nil {
		t.Errorf("VSIF.Seek(0, 0) error is not nil: %s", err)
	}
	if offset != 0 {
		t.Errorf("VSIF.Seek(0, 0) == %d, expected 0", offset)
	}

	// Test Read
	out := make([]byte, 20, 20)
	c, err = fp.Read(out)
	if err == nil {
		t.Errorf("VSIF.Read() error is nil, expected EOF")
	}
	if err != io.EOF {
		t.Errorf("VSIF.Read() error is not EOF: %s", err)
	}
	if c != len(in) {
		t.Errorf("VSIF.Read() bytes read == %d, expected %d", c, len(in))
	}
	if string(out[:c]) != string(in) {
		t.Errorf("VSIF.Read() == %q, not %q", string(out[:c]), string(in))
	}

	// Test Truncate
	expected := in[:7]
	size := len(expected)
	err = fp.Truncate(int64(size))
	if err != nil {
		t.Errorf("VSIF.Truncate(%v) error is not nil: %s", size, err)
	}
	offset, err = fp.Seek(0, 0)
	if err != nil {
		t.Errorf("VSIF.Seek(0, 0) error is not nil: %s", err)
	}
	if offset != 0 {
		t.Errorf("VSIF.Seek(0, 0) == %d, expected 0", offset)
	}
	c, err = fp.Read(out)
	if err == nil {
		t.Errorf("VSIF.Read() error is nil, expected EOF")
	}
	if err != io.EOF {
		t.Errorf("VSIF.Read() error is not EOF: %s", err)
	}
	if c != size {
		t.Errorf("VSIF.Read() bytes read == %d, expected %d", c, size)
	}
	if string(out[:c]) != string(expected) {
		t.Errorf("VSIF.Read() == %q, not %q", string(out[:c]), string(expected))
	}

	// Test Close
	err = fp.Close()
	if err != nil {
		t.Errorf("VSIF.Close() error is not nil: %s", err)
	}
}

func TestVSIFReadWriter(t *testing.T) {
	// Open the CSV
	fname := "/vsimem/test.csv"
	fp, err := VSIFOpen(fname, "rw")
	if err != nil {
		t.Errorf("VSIFOpen(%v) error is not nil: %s", fname, err)
	}

	// The records to add.
	input := [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}

	w := csv.NewWriter(fp)
	w.WriteAll(input) // calls Flush internally

	if err = w.Error(); err != nil {
		t.Errorf("Writing to CSV failed: %s", err)
	}

	// Go back to the beginning of the file.
	var offset int64
	offset, err = fp.Seek(0, 0)
	if err != nil {
		t.Errorf("VSIF.Seek(0, 0) error is not nil: %s", err)
	}
	if offset != 0 {
		t.Errorf("VSIF.Seek(0, 0) == %d, expected 0", offset)
	}

	var output [][]string
	r := csv.NewReader(fp)
	output, err = r.ReadAll()
	if err != nil {
		t.Errorf("Reading from CSV vailed: %s", err)
	}

	if !reflect.DeepEqual(input, output) {
		t.Errorf("Input records != ouput")
	}
}
