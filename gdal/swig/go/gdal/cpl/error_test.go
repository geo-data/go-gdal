package cpl_test

import (
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCustomHandlers(t *testing.T) {
	var globalCount, localCount int32
	globalHandler := cpl.ErrorHandler(func(e cpl.Error, data interface{}) {
		atomic.AddInt32(&globalCount, 1)
		//log.Printf("Global error handler: error type %d, number %d: %s", e.ErrorType(), e.ErrorNum(), e.Error())
	})
	cpl.SetErrorHandler(globalHandler, nil)

	err := cpl.NewError(constant.CE_Warning, constant.CPLE_AppDefined, "first global error")
	cpl.SetError(err)

	var wg sync.WaitGroup
	workers := 3
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			runtime.LockOSThread() // This is released when the goroutine exits.
			defer wg.Done()

			localHandler := cpl.ErrorHandler(func(e cpl.Error, data interface{}) {
				atomic.AddInt32(&localCount, 1)
				//log.Printf("Worker %d error handler: error type %d, number %d: %s", id, e.ErrorType(), e.ErrorNum(), e.Error())
			})
			cpl.PushErrorHandler(localHandler, nil)

			msg := fmt.Sprintf("worker %d local error", id)
			lerr := cpl.NewError(constant.CE_Failure, constant.CPLE_AppDefined, msg)
			cpl.SetError(lerr)

			cpl.PopErrorHandler()

			msg = fmt.Sprintf("worker %d global error", id)
			lerr = cpl.NewError(constant.CE_Warning, constant.CPLE_AppDefined, msg)
			cpl.SetError(lerr)
		}(i)
	}

	wg.Wait()

	err = cpl.NewError(constant.CE_Warning, constant.CPLE_AppDefined, "second global error")
	cpl.SetError(err)

	globalExpected := int32(workers + 2)
	if globalCount != globalExpected {
		t.Errorf("ErrorHandler: expected %d calls to global handler, got %d", globalExpected, globalCount)
	}

	if localCount != int32(workers) {
		t.Errorf("ErrorHandler: expected %d calls to local handler, got %d", workers, localCount)
	}
}

// TestErrorTrap tests the error trapping idiom used throughout the package.
func TestErrorTrap(t *testing.T) {

	// Trigger and trap an error.
	triggerError := func() (err error) {
		defer cpl.ErrorTrap()(&err)

		// Generate a new error.
		e := cpl.NewError(constant.CE_Warning, constant.CPLE_AppDefined, "test error")
		cpl.SetError(e)
		return
	}

	// Have we trapped an error?
	err := triggerError()
	if err == nil {
		t.Fatal("ErrorTrap(): error is nil")
	}

	// Is the error of the expected type?
	cerr, ok := err.(cpl.Error)
	if !ok {
		t.Fatalf("ErrorTrap(): expected error type cpl.Error, got %T", err)
	}

	// Does the error have the expected number?
	if cerr.ErrorNum() != constant.CPLE_AppDefined {
		t.Errorf("ErrorTrap(): error number == %v, expected %v", cerr.ErrorNum(), constant.CPLE_AppDefined)
	}

	// Does the error have the expected type?
	if cerr.ErrorType() != constant.CE_Warning {
		t.Errorf("ErrorTrap(): error type == %v, expected %v", cerr.ErrorNum(), constant.CE_Warning)
	}

	// Does the error have the expected message?
	if msg := cerr.Error(); msg != "test error" {
		t.Errorf("ErrorTrap(): error message == %v, expected \"test error\"", msg)
	}

	// Do we have a previous error?
	last := cpl.LastError()
	if last == nil {
		t.Fatal("LastError(): error is nil")
	}

	// Is the previous error the same as the triggered error?
	if last.Error() != err.Error() {
		t.Errorf("LastError() == %v, expected %v", last.Error(), err.Error())
	}
}
