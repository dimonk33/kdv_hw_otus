package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 || m <= 0 {
		return errors.New("wrong initial parameters")
	}
	var errCount int32 = 0
	chTask := make(chan Task)
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range chTask {
				if task() != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}

	for _, t := range tasks {

		if atomic.LoadInt32(&errCount) >= int32(m) {
			break
		}
		chTask <- t
	}
	close(chTask)

	wg.Wait()

	if errCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
