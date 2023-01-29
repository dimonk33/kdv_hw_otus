package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	curTaskIdx := 0
	errCount := 0
	muTask := sync.Mutex{}
	muError := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(_curTaskIdx, _errCount *int) {
			var task Task

			for {
				task = nil
				muTask.Lock()
				if *_curTaskIdx < len(tasks) {
					task = tasks[*_curTaskIdx]
					*_curTaskIdx++
				}
				muTask.Unlock()
				if task == nil {
					break
				}

				if task() != nil {
					muError.Lock()
					*_errCount++
					muError.Unlock()
				}

				muError.Lock()
				count := *_errCount
				muError.Unlock()

				if count >= m {
					break
				}
			}
			wg.Done()
		}(&curTaskIdx, &errCount)
	}

	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
