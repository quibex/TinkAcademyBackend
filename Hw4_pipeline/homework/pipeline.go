package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {

	for _, stage := range stages {
		//соединяем выход первого этапа с входом следующего и т.д.
		in = stage(in)
	}
	out := make(chan any)
	go func() { //чтобы контролировать завершение
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case elem, ok := <-in:
				if !ok {
					return
				}
				out <- elem
			}
		}
	}()

	return out
}
