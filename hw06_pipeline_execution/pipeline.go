package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	mu := sync.RWMutex{}
	out := make(Bi)

	id := 0
	order := 0
	for val := range in {
		go stageThread(id, &order, &mu, stages, val, out, done)
		id++
	}

	go func(_order *int) {
		defer close(out)
		exit := false
		for !exit {
			select {
			case _, ok := <-done:
				if !ok {
					return
				}
			default:
				mu.RLock()
				exit = *_order >= id
				mu.RUnlock()
			}
		}
	}(&order)

	return out
}

func stageThread(
	id int,
	order *int,
	mu *sync.RWMutex,
	stages []Stage,
	val interface{},
	out chan<- interface{},
	done In,
) {
	ch := make(Bi, 1)
	ch <- val
	close(ch)

	var _out Out = ch
	for _, st := range stages {
		select {
		case _, ok := <-done:
			if !ok {
				return
			}
		default:
			_out = st(_out)
		}
	}
	val = <-_out

	exit := false
	for !exit {
		select {
		case _, ok := <-done:
			if !ok {
				return
			}
		default:
			mu.RLock()
			exit = *order >= id
			mu.RUnlock()
		}
	}

	exit = false
	for !exit {
		select {
		case _, ok := <-done:
			if !ok {
				return
			}
		case out <- val:
			mu.Lock()
			*order++
			mu.Unlock()
			return
		}
	}
}
