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

type TaskFunc func(context.Context)

type Pool interface {
	Go(task TaskFunc)
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
