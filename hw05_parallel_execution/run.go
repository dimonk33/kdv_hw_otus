package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 || m <= 0 {
		return errors.New("wrong initial parameters")
	}
	errCount := 0
	muError := sync.Mutex{}
	chTask := make(chan Task)
	wg := sync.WaitGroup{}
	wg.Add(n + 1)

	for i := 0; i < n; i++ {
		go func(chT <-chan Task, _errCount *int) {
			defer wg.Done()
			for task := range chT {
				if task() != nil {
					muError.Lock()
					*_errCount++
					muError.Unlock()
				}
			}
		}(chTask, &errCount)
	}

	go func(chT chan<- Task, _errCount *int) {
		defer close(chT)
		defer wg.Done()

		for _, t := range tasks {
			wait := true
			for wait {
				muError.Lock()
				exit := *_errCount >= m
				muError.Unlock()
				if exit {
					return
				}

				select {
				case chT <- t:
					wait = false
				default:
				}
			}
		}
	}(chTask, &errCount)

	wg.Wait()

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
