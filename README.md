# Rocket (in development)
A rest API server library focused in speed and simplicity

## Install
```sh
$ go get github.com/erichnascimento/rocket
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

	// Add logger middleware
	s.Use(logger.NewLogger())

	// Add router middleware prefixed for "/myapi/v2"
	r := router.NewRouter("/myapi/v2")

	// add a simple route
	r.Add("GET", "/test", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "Welcome!\n")
	})

	// add a route
	r.Add("GET", "/users", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "List users!\n")
	})

	// add route with param
	r.Add("GET", "/users/:userId", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "Listing user %s!\n", ctx.GetParam("userId"))
	})

	// add route and subroute
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
2015/09/08 16:10:44 GET /myapi/v2/test 200 0ms - 9B
2015/09/08 16:10:54 GET /myapi/v2/users 200 0ms - 12B
2015/09/08 16:11:00 GET /myapi/v2/users/123 200 0ms - 18B
2015/09/08 16:11:08 GET /myapi/v2/users/123/sales/456 200 0ms - 31B

```
## Perform requests

```sh
$ curl http://127.0.0.1:3000/myapi/v2/test
Welcome!

$ curl http://127.0.0.1:3000/myapi/v2/users
List users!

$ curl http://127.0.0.1:3000/myapi/v2/users/123
Listing user 123!

$ curl http://127.0.0.1:3000/myapi/v2/users/123/sales/456
Listing user 123 and sale 456!

```
