package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tdv/go-pool/pool"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	p := pool.New(ctx, 5)
	defer p.Stop()

	for i := 0; i < 20; i++ {
		i := i
		p.Go(func(context.Context) {
			fmt.Println(i)
			time.Sleep(time.Second * 8)
		})
	}

	time.Sleep(time.Second * 20)
}
