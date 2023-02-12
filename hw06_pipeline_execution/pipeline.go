package hw06pipelineexecution

import (
	"fmt"
	"time"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var from = in
	start := time.Now()
	for i, stage := range stages {
		inCh := make(Bi)

		go func(id int, in_ In, out_ Bi, stage_ Stage) {
			defer close(out_)
			fmt.Printf("stage %d start %v\r\n", id, time.Since(start))
			stCh := stage_(in_)
			for val := range stCh {
				select {
				case <-done:
					return
				default:
					out_ <- val
				}
			}
			fmt.Printf("stage %d finished %v\r\n", id, time.Since(start))
		}(i, from, inCh, stage)

		from = inCh
	}

	return from
}
