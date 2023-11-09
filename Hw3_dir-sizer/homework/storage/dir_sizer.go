package storage

import (
	"context"
	"sync"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	chRes := make(chan Result, a.maxWorkersCount)
	chErr := make(chan error)
	wg := &sync.WaitGroup{}
	res := Result{}

	dirs, files, err := d.Ls(ctx)
	if err != nil {
		return Result{}, err
	}

	for i := range dirs {
		wg.Add(1)
		go func(d Dir) {
			defer wg.Done()
			r, err := a.Size(ctx, d)
			if err != nil {
				chErr <- err
				return
			}
			chRes <- r
		}(dirs[i])
	}

	for i := range files {
		sizeFile, err := files[i].Stat(ctx)
		if err != nil {
			return Result{}, err
		}
		res.Size += sizeFile
		res.Count++
	}

	go func() {
		wg.Wait()
		close(chRes)
		close(chErr)
	}()

	for {
		select {
		case err := <-chErr:
			if err != nil {
				return Result{}, err
			}
		case resI, ok := <-chRes:
			if !ok {
				return res, nil
			}
			res.Size += resI.Size
			res.Count += resI.Count
		}
	}
}
