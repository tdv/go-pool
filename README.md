# go-pool
go-pool - is a library goroutine pools conceptually close to the thread pool in other languages.

The common gist is to limit concurrency and execute only limited number tasks at the time.

Goroutines are well done and appropriate for many tasks but sometime it is an issue. The run immediately after calling 'Go'  might reduce available of the other resources like DB for instance.

There is an issue
```go
func doSomething(n int, connectionString string) {
	for i := 0; i < n; i++ {
		go func() {
			db, err := sql.Open(connectionString)
			// Something else here ...
			defer db.Close()
			// Something else here ...
		}()
	}
}

func main() {
	//Your DB has died. May be...
	doSomething(100500, "Does not metter")
	// ...
}
```

A proposed solution

```go
import (
	"context"

	"github.com/tdv/go-pool/pool"
)

func doSomething(p pool.Pool, n int, connectionString string) {
	for i := 0; i < n; i++ {
		p.Go(func() {
			db, err := sql.Open(connectionString)
			// Something else here ...
			defer db.Close()
			// Something else here ...
		})
	}
}

func main() {
	ctx := context.Background()
	p := pool.New(ctx, 5)
	defer p.Stop()

	// There is no any problem.
	// All will go through 5 DB connection.
	// DB will be healthy and your mind too.
	doSomething(p, 100500, "Does not metter")
	// ...
}
```

# Installation

```bash
go get github.com/tdv/go-pool/pool
```



# Examples

## example1
[Source code](https://github.com/tdv/go-pool/tree/main/examples/example1)
**Description**  
The example demonstrates the normal work of the limited gorutine pool.
```go
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
```

## example2
[Source code](https://github.com/tdv/go-pool/tree/main/examples/example2)
**Description**  
The example demonstrates the cancelation of the work. In spite of expectation to see all 20 lines in terminal there will be only 5 and the others will be canceled.
```go
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
			time.Sleep(time.Second * 7)
		})
	}

	time.Sleep(time.Second * 10)
}
```

## example3
[Source code](https://github.com/tdv/go-pool/tree/main/examples/example3)
**Description**  
The example demonstrates how to decline all tasks in queue manually.
```go
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
```

# Wrapping up
Goroutines in Go are really useful, but sometimes we need to solve some design issues. The proposed library is one of many solutions.
