package main

import (
	"net/http"
	"strconv"

	"log"

	"github.com/erichnascimento/rocket/middleware"
	"github.com/erichnascimento/rocket/server"
	"github.com/erichnascimento/rocket/server/response"
)

func main() {
	s := server.NewServer()

	// Use a Logger middleware for logging
	s.Use(middleware.NewLogger())

	// Create a new router middleware for API `my_api`, version 2
	r := middleware.NewRouter("/my_api/v2")
	r.Get(`/users`, listUsersHandler)
	r.Get(`/users/:id`, getUserHandler)

	// use the router
	s.Use(r)

	// Start listening and serving
	err := s.ListenAndServe(":2000")
	log.Print(err)
}

func listUsersHandler(rw http.ResponseWriter, _ *http.Request) {
	response.SendJSON(rw, users, http.StatusOK)
}

func getUserHandler(rw http.ResponseWriter, req *http.Request) {
	pID := req.Context().Value(`id`).(string)
	id, err := strconv.Atoi(pID)
	if err != nil {
		response.SendJSON(rw, "id arguments should be a integer", http.StatusBadRequest)
		return
	}
	if id < 1 || id > 2 {
		response.SendJSON(rw, "User not found", http.StatusNotFound)
		return
	}

	response.SendJSON(rw, users[id-1], http.StatusOK)
}

var users = []user{
	user{1, "Jacob"},
	user{2, "Dudu"},
}

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
