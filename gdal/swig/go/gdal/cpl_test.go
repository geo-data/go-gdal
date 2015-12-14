package gdal_test

import (
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"os"
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	var empty string
	var dirtests = []struct {
		dir     string
		err     bool
		compare bool // Compare against golang stdlib equivalent.
	}{
		{empty, false, false},
		{"", false, false},
		{".", false, true},
		{"..", false, true},
		{"../", false, true},
		{"../.", false, true},
		{"oops", false, false},        // does not exist
		{"cpl_test.go", false, false}, // a file
		{"/vsimem/oops", false, false},
	}

	for _, dt := range dirtests {
		dir := dt.dir

		var names, expected []string
		var err error

		names, err = gdal.ReadDir(dir)
		if dt.err {
			if err == nil {
				t.Errorf("ReadDir(%q): expected error but none found", dir)
			}
			continue
		}

		if !dt.err && err != nil {
			t.Errorf("ReadDir(%q): unexpected error, %v", dir, err)
			continue
		}

		if !dt.compare {
			continue
		}

		f, e := os.Open(dir)
		if err != nil {
			t.Errorf("os.Open(%q): unexpected error, %v", dir, e)
			continue
		}

		expected, err = f.Readdirnames(-1)
		if err != nil {
			t.Errorf("f.Readdirnames() for %v: unexpected error, %v", dir, e)
			continue
		}

	Check:
		for _, name := range names {
			for i, e := range expected {
				if e != name {
					continue
				}
				expected = append(expected[:i], expected[i+1:]...) // delete the expected name
				if len(expected) == 0 {
					break Check
				}
			}
		}

		if len(expected) > 0 {
			t.Errorf("ReadDir(%v): expected files not found: %v", dir, expected)
		}
	}
}

func TestParseCommandLine(t *testing.T) {
	var empty string
	var parsetests = []struct {
		opts   string
		expect []string
	}{
		{empty, nil},
		{"", nil},
		{"foo", []string{"foo"}},
		{"foo bar", []string{"foo", "bar"}},
		{"--test 1 -t \"foo bar\"", []string{"--test", "1", "-t", "foo bar"}},
	}

	for _, pt := range parsetests {
		opts := pt.opts
		expect := pt.expect

		result := gdal.ParseCommandLine(opts)

		if !reflect.DeepEqual(result, expect) {
			t.Errorf("ParseCommandLine(%q) == %q, does not match %q", opts, result, expect)
		}
	}
}
