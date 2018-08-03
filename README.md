# Rocket (WIP)
A HTTP rest API server library focused on performance and simplicity

## Install
```sh
$ go get github.com/erichnascimento/rocket
```

## Usage
```go

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


```

## Running the simple-server example

```sh
$ go run examples/simple-server.go
2018/08/02 22:41:12 GET /my_api/v2/users 200 0ms - 49 B
2018/08/02 22:42:25 GET /my_api/v2/users/1 200 0ms - 24 B
2018/08/02 22:44:06 GET /my_api/v2/users/3 404 0ms - 17 B

```
## Perform requests

```sh
$ curl -XGET "http://localhost:2000/my_api/v2/users" -v
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 03 Aug 2018 01:41:12 GMT
< Content-Length: 49
<
[{"id":1,"name":"Jacob"},{"id":2,"name":"Dudu"}]

$ curl -XGET "http://localhost:2000/my_api/v2/users/1" -v
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 03 Aug 2018 01:42:25 GMT
< Content-Length: 24
<
{"id":1,"name":"Jacob"}

$ curl -XGET "http://localhost:2000/my_api/v2/users/3" -v
< HTTP/1.1 404 Not Found
< Content-Type: application/json
< Date: Fri, 03 Aug 2018 01:44:06 GMT
< Content-Length: 17
<
"User not found"

```
