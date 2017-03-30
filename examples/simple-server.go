package main

import (
	"net/http"
	"fmt"

	"github.com/erichnascimento/rocket/server"
	"github.com/erichnascimento/rocket/middleware"
	"github.com/erichnascimento/rocket/middleware/logger"
)

func main() {
	s := server.NewRocket()

	// Use a Logger middleware for logging
	s.Use(logger.NewLogger())

	// Create a new router middleware for API `my_api`, version 2
	r := middleware.NewRouter("/my_api/v2")
	r.Get(`/users`, listUsersHandler)
	r.Get(`/users/:id`, getUserHandler)

	// use the router
	s.Use(r)

	// Start listening and serving
	s.ListenAndServe(":2000")
}

var users = map[interface{}]string{
	`1`:"Jacob",
	`2`:"Dudu",
}

func listUsersHandler(rw http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(rw, `[`)
	for id, name := range users {
		fmt.Fprintf(rw, `{"id": %d, "name": "%s"},`, id, name)
	}
	fmt.Fprint(rw, `]`)
}

func getUserHandler(rw http.ResponseWriter, req *http.Request) {
	id := req.Context().Value(`id`)
	user := users[id]
	fmt.Fprint(rw, user)
}
