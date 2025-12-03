package future

import "sync"

// run all functions and wait for them all finish
func Wait(funcs ...func()) {
	Run(funcs...).Wait()
}

// run all functions and don't wait for finish
func Run(funcs ...func()) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(len(funcs))

	for _, f := range funcs {
		go func(fn func()) {
			defer wg.Done()
			fn()
		}(f)
	}

	return &wg
}
