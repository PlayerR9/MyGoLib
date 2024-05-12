package Runner

import (
	"sync"

	mext "github.com/PlayerR9/MyGoLib/Utility/MapExt"
)

// WaitAll is a function that waits for all Go routines in the batch to finish
// and returns a slice of errors that represent the error statuses of the Go routines.
//
// Parameters:
//   - batch: A slice of pointers to the GRHandler instances that handle the Go routines.
//
// Returns:
//   - []error: A slice of errors that represent the error statuses of the Go routines.
func WaitAll(batch map[string]*HandlerSimple) []error {
	// 1. Filter out nil GRHandler instances.
	batch = mext.FilterNilValues(batch)
	if len(batch) == 0 {
		return nil
	}

	// 2. Initialize all error channels.
	errChans := make([]<-chan error, 0, len(batch))

	for _, h := range batch {
		errChans = append(errChans, h.GetErrChannel())
	}

	// 3. Wait for all Go routines to finish.
	errs := make([]error, len(batch))

	var wg sync.WaitGroup

	wg.Add(len(batch))

	for i, errChan := range errChans {
		go func(i int, errChan <-chan error) {
			defer wg.Done()

			for err := range errChan {
				errs[i] = err
			}
		}(i, errChan)
	}

	wg.Wait()

	return errs
}
