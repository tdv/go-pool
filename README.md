# go-pool
go-pool - is a little library goroutine pools conceptually close to the thread pool in other languages.

The common gist is to limit concurrency and execute only a limited number of the tasks at the time.

Goroutines are well done and appropriate for many tasks but sometimes it is an issue. The run immediately after calling 'ao'  might reduce the availability of the other resources like DB for instance.



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
  // Your DB has died. May be...
  doSomething(100500, "Does not metter")
  // ...
}
```

A proposed solution

```go
import (
  "context"

  "github.com/tdv/go-pool"
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

  // There is no problem.
  // All will go through 5 DB connections.
  // DB will be healthy and your mind too.
  doSomething(p, 100500, "Does not metter")
  // ...
}
```

# Installation
```bash
go get -u github.com/tdv/go-pool
```

# Examples
## Normal work
[example1](https://github.com/tdv/go-pool/tree/main/examples/example1)  
**Description**  
The example demonstrates the normal work of the limited goroutine pool.  
```go
package main

import (
  "context"
  "fmt"
  "sync"
  "time"

  "github.com/tdv/go-pool"
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
## Cancelation through the context
[example2](https://github.com/tdv/go-pool/tree/main/examples/example2)  
**Description**  
The example demonstrates the cancellation of the work. In spite of the expectation to see all 20 lines in the terminal there will be only 10 (2 first pieces / portions) and the others will be canceled.  
```go
package main

import (
  "context"
  "fmt"
  "time"

  "github.com/tdv/go-pool"
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
## Manual pool closing
[example3](https://github.com/tdv/go-pool/tree/main/examples/example3)  
**Description**  
The example demonstrates how to decline all tasks in the queue manually. In spite of the expectation to see all 20 lines in the terminal there will be only 5 (first 5 tasks have been done before the pool is closed) and the others will be canceled.  
```go
package main

import (
  "context"
  "fmt"
  "time"

  "github.com/tdv/go-pool"
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
