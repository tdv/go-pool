package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tdv/go-pool/pool"
)

func main() {
	ctx := context.Background()

	p := pool.New(ctx, 5)

	for i := 0; i < 20; i++ {
		i := i

		p.Go(func(context.Context) {
			fmt.Println(i)
			time.Sleep(time.Second * 2)
		})
	}

	time.Sleep(time.Second * 1)
	p.Stop()
}
