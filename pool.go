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

// A TaskFunc is an func for execution, which will be executed
// with respective concurrency limits.
// The function gets a context, via which you can make
// graceful cancellation within you function.
type TaskFunc func(context.Context)

// Conceptually, the Pool looks like a thread pool in the other languages,
// but implemented on the goroutines.
// The gist is to limit the number of the concurrent tasks
// at the same time of execution.
type Pool interface {
	// Add a task into the pool.
	//
	// The task will executed in the normal flow
	// or canceled if the pool is stopped (or context is canceled)
	Go(task TaskFunc)

	// Stopping executing any tasks.
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

func (this *poolImpl) Go(task TaskFunc) {
	this.wg.Add(1)

	ctx, cancel := context.WithCancel(this.ctx)
	go func() {
		defer this.wg.Done()
		defer cancel()

		select {
		case this.slots <- struct{}{}:
			if this.ctx.Err() == nil {
				task(ctx)
				<-this.slots
			}
		case <-this.ctx.Done():
			return
		}
	}()
}

func (this *poolImpl) Stop() {
	this.cancel()
	this.wg.Wait()
}

// Create a new pool.
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
