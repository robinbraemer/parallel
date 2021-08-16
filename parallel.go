package parallel

import (
	"go.uber.org/multierr"
	"sync"
)

// Do is the interface for an executable function/task.
type Do interface{ Do() error }

// Fn implements Do interface.
type Fn func() error

func (fn Fn) Do() error { return fn() }

// Parallel returns a Do that when run executes all passed Do's in parallel
// and blocks until all are done returning nil, one error or a multi-error.
func Parallel(do ...Do) Do {
	return Fn(func() error {
		switch len(do) {
		case 0:
			return nil
		case 1:
			return do[0].Do()
		}
		var (
			wg      sync.WaitGroup
			errs    []error
			errChan = make(chan error)
		)
		wg.Add(len(do))
		go func() { wg.Wait(); close(errChan) }()
		for _, fn := range do {
			fn := fn
			go func() {
				defer wg.Done()
				if err := fn.Do(); err != nil {
					errChan <- err
				}
			}()
		}
		for err := range errChan {
			errs = append(errs, err)
		}
		return multierr.Combine(errs...)
	})
}

// Ordered returns a Do that when run executes all passed Do's in order
// and blocks until all are done returning nil or the first error occurred.
func Ordered(do ...Do) Do {
	return Fn(func() error {
		for _, fn := range do {
			if err := fn.Do(); err != nil {
				return err
			}
		}
		return nil
	})
}
