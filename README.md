# Rocket (in development)
A rest API server library focused in speed and simplicity

## Install
```sh
$ go get github.com/erich
```

## Usage
```go

package main

import (
	"fmt"

	"github.com/erichnascimento/rocket/middleware/logger"
	"github.com/erichnascimento/rocket/middleware/router"
	"github.com/erichnascimento/rocket/server"
)

func main() {
	s := server.New("0.0.0.0:3000")

	s.Use(logger.NewLogger())

	r := router.NewRouter()
	r.Add("GET", "/test", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "Welcome!\n")
	})

	r.Add("GET", "/users", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "List users!\n")
	})

	r.Add("GET", "/users/:userId", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "Listing user %s!\n", ctx.GetParam("userId"))
	})

	r.Add("GET", "/users/:userId/sales/:saleId", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "Listing user %s and sale %s!\n", ctx.GetParam("userId"), ctx.GetParam("saleId"))
	})

	s.Use(r)
	s.Serve()
}


```

## Running the simple-server example

```sh
$ go run examples/simple-server.go
2015/09/04 18:01:02 GET /test 200 0ms - 9B
2015/09/04 18:01:11 GET /users 200 0ms - 12B
2015/09/04 18:01:15 GET /users/123 200 0ms - 18B
2015/09/04 18:01:18 GET /users/123/sales/456 200 0ms - 31B

```
## Perform requests

```sh
$ curl http://127.0.0.1:3000/test
Welcome!

$ curl http://127.0.0.1:3000/users
List users!

$ curl http://127.0.0.1:3000/users/123
Listing user 123!

$ curl http://127.0.0.1:3000/users/123/sales/456
Listing user 123 and sale 456!

```
