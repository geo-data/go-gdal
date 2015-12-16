package main

import (
	"fmt"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/constant"
	"github.com/geo-data/go-gdal/gdal/swig/go/gdal/cpl"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	var count int32 = 0
	globalHandler := cpl.ErrorHandler(func(e cpl.Error, data interface{}) {
		atomic.AddInt32(&count, 1)
		log.Printf("Global error handler: error type %d, number %d: %s", e.ErrorType(), e.ErrorNum(), e.Error())
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
				atomic.AddInt32(&count, 1)
				log.Printf("Worker %d error handler: error type %d, number %d: %s", id, e.ErrorType(), e.ErrorNum(), e.Error())
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

	log.Printf("Finished: %d errors logged\n", count)
}
