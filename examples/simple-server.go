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

	// Add router middleware
	r := router.NewRouter()
	r.Add("GET", "/test", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "Welcome!\n")
	})

	// add a route
	r.Add("GET", "/users", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "List users!\n")
	})

	// add route
	r.Add("GET", "/users/:userId", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "Listing user %s!\n", ctx.GetParam("userId"))
	})

	// add route
	r.Add("GET", "/users/:userId/sales/:saleId", func(ctx *router.Context) {
		fmt.Fprintf(ctx, "Listing user %s and sale %s!\n", ctx.GetParam("userId"), ctx.GetParam("saleId"))
	})

	s.Use(r)
	s.Serve()
}
