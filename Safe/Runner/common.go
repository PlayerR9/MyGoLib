package Runner

import (
	"sync"

	mext "github.com/PlayerR9/MyGoLib/CustomData/OrderedMap"
)

// WaitAll is a function that waits for all Go routines in the batch to finish
// and returns a slice of errors that represent the error statuses of the Go routines.
//
// Parameters:
//   - batch: A slice of pointers to the GRHandler instances that handle the Go routines.
//
// Returns:
//   - []error: A slice of errors that represent the error statuses of the Go routines.
func WaitAll(batch map[string]*HandlerSimple) map[string]error {
	// 1. Filter out nil GRHandler instances.
	batch = mext.FilterNilValues(batch)
	if len(batch) == 0 {
		return nil
	}

	// 2. Initialize all error channels.
	errChans := make(map[string]<-chan error)

	for k, h := range batch {
		errChans[k] = h.GetErrChannel()
	}

	// 3. Wait for all Go routines to finish.
	errs := make(map[string]error)

	for k := range batch {
		errs[k] = nil
	}

	var wg sync.WaitGroup

	wg.Add(len(batch))

	for k, errChan := range errChans {
		go func(k string, errChan <-chan error) {
			defer wg.Done()

			for err := range errChan {
				errs[k] = err
			}
		}(k, errChan)
	}

	wg.Wait()

	return errs
}
