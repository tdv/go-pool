// Package:	go-pool (https://github.com/tdv/go-pool)
// Created:	09.2022
// Copyright 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Distributed under the MIT License
// (See accompanying file LICENSE)

package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tdv/go-pool/pool"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	p := pool.New(ctx, 5)
	defer p.Stop()

	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		i := i
		wg.Add(1)

		p.Go(func(context.Context) {
			defer wg.Done()
			fmt.Println(i)
			time.Sleep(time.Second * 2)
		})
	}

	time.Sleep(time.Second)
	wg.Wait()
}
