package work

import (
	"context"
	"errors"
	"sync"
)

type Executor interface {
	Execute() error
	OnError(error)
}

type Pool struct {
	maxWorkers int
	tasks      chan Executor
	start      sync.Once
	stop       sync.Once
	quit       chan struct{}
}

func NewPool(max int, taskChanSize int) (*Pool, error) {
	if max <= 0 {
		return nil, errors.New("worker cannot be less, or equal to zero")
	}

	if taskChanSize < 0 {
		return nil, errors.New("channel size cannot be a negative value")
	}

	return &Pool{
		maxWorkers: max,
		tasks:      make(chan Executor, taskChanSize),
		start:      sync.Once{},
		stop:       sync.Once{},
		quit:       make(chan struct{}),
	}, nil
}

func (p *Pool) Start(ctx context.Context) {
	p.start.Do(func() {
		p.startWorker(ctx)
	})
}

func (p *Pool) Stop() {
	p.stop.Do(func() {
		close(p.quit)
	})
}

func (p *Pool) AddTask(t Executor) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}

func (p *Pool) startWorker(ctx context.Context) {
	for i := 0; i < p.maxWorkers; i++ {
		go func(workerNum int) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-p.quit:
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}
					if err := task.Execute(); err != nil {
						task.OnError(err)
					}
				}
			}
		}(i)
	}
}
