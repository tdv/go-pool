// Package:	go-pool (https://github.com/tdv/go-pool)
// Created:	09.2022
// Copyright 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Distributed under the MIT License
// (See accompanying file LICENSE)

package pool

import (
	"context"
	"sync"
)

// A TaskFunc is a func for execution, which will be executed
// with respective concurrency limits.
// The function gets a context, via which you can make
// graceful cancellation within you function.
type TaskFunc func(context.Context)

// Pool - conceptually, the Pool looks like a thread pool in the other languages,
// but implemented on the goroutines.
// The gist is to limit the number of the concurrent tasks
// at the same time of execution.
type Pool interface {
	// Go - adds a task into the pool.
	//
	// The task will execute in the normal flow
	// or canceled if the pool is stopped (or context is canceled)
	Go(task TaskFunc)

	// Stop - stops all tasks.
	//
	// All tasks in the queue to execute will be canceled.
	// The cancellation will propagate via the Context through all the tasks.
	// Call the method only once for instance.
	Stop()
}

type poolImpl struct {
	ctx    context.Context
	cancel context.CancelFunc
	slots  chan struct{}
	wg     sync.WaitGroup
}

func (s *poolImpl) Go(task TaskFunc) {
	s.wg.Add(1)

	ctx, cancel := context.WithCancel(s.ctx)
	go func() {
		defer s.wg.Done()
		defer cancel()

		select {
		case s.slots <- struct{}{}:
			if s.ctx.Err() == nil {
				task(ctx)
				<-s.slots
			}
		case <-s.ctx.Done():
			return
		}
	}()
}

func (s *poolImpl) Stop() {
	s.cancel()
	s.wg.Wait()
}

// New - creates a new pool.
//
// The function gets context and a limit of
// a number of the concurrent tasks.
func New(ctx context.Context, capacity uint) Pool {
	if capacity < 1 {
		capacity = 1
	}

	ctx, cancel := context.WithCancel(ctx)

	return &poolImpl{
		ctx:    ctx,
		cancel: cancel,
		slots:  make(chan struct{}, capacity),
	}
}
