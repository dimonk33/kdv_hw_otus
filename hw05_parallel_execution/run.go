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
		go func(_m, _curTaskIdx, _errCount *int, _tasks []Task) {
			for *_errCount < *_m && *_curTaskIdx < len(_tasks) {
				muTask.Lock()
				task := _tasks[*_curTaskIdx]
				*_curTaskIdx++
				muTask.Unlock()
				if task() != nil {
					muError.Lock()
					*_errCount++
					muError.Unlock()
				}
			}
			wg.Done()
		}(&m, &curTaskIdx, &errCount, tasks)
	}

	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
