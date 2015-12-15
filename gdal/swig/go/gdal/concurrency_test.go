package gdal_test

import (
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// TestConcurrency() tests whether errors are properly synchronised between
// goroutines and gdal.  It does this by concurrently opening 250,000 unique non
// existent datasets in separate goroutines: each goroutine should generate an
// unique error.  Each error is passed via a channel to a corresponding
// goroutine which checks whether the error is correctly associated with the
// expected name.
func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	count := 250000
	wg.Add(count)

	for i := 0; i < count; i++ {
		done := make(chan error, 1)
		name := fmt.Sprintf("non-existent-%d", i)

		// Open the non existent dataset, generating an error.
		go func(name string, done chan error) {
			runtime.LockOSThread() // This is released when the goroutine exits.
			_, err := gdal.Open(name, constant.OF_READONLY|constant.OF_VERBOSE_ERROR)
			if err == nil {
				err = fmt.Errorf("Open(%q) error == nil", name)
				defer t.Fatal(err)
			}
			done <- err
		}(name, done)

		// Check the error against the name.
		go func(name string, done chan error) {
			defer wg.Done()
			err := <-done
			if !strings.Contains(err.Error(), name) {
				t.Errorf("Open(%q) name not in error: %s\n", name, err.Error())
			}
		}(name, done)
	}

	wg.Wait()
}
