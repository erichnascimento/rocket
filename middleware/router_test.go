package middleware

import (
	"testing"
	"net/http"

	"github.com/erichnascimento/rocket"

	"github.com/bmizerany/assert"
)

func TestGetPath(t *testing.T) {
	paths := map[string] string{
		"/user": "/user",
		"/user?": "/user",
		"/user?p=123": "/user",
		"/user#app/list": "/user",
	}

	for k, v := range paths {
		assert.Equal(t, getPath(k), v)
	}
}

func TestEmptyRoute(t *testing.T) {
	err, r := newRoute("", nil)

	assert.Equal(t, err, nil)
	assert.NotEqual(t, r, nil)
	assert.Equal(t, r.route, "")
	assert.Equal(t, r.compiledRoute, ".")
	assert.Equal(t, len(r.resources), 0)
	assert.Equal(t, len(r.params), 0)
}

func TestFullRoute(t *testing.T) {
	_, r := newRoute("/users/:userId/sales/:saleId?q=123", nil)

	assert.Equal(t, r.compiledRoute, "./users/$/sales/$")

	assert.Equal(t, len(r.resources), 2)
	assert.Equal(t, r.resources["users"], true)
	assert.Equal(t, r.resources["sales"], true)

	assert.Equal(t, len(r.params), 2)
	assert.Equal(t, r.params[0], "userId")
	assert.Equal(t, r.params[1], "saleId")
}

func TestEmptyRequest(t *testing.T) {
	err, r := newRequest("", nil)

	assert.Equal(t, err, nil)
	assert.NotEqual(t, r, nil)
	assert.Equal(t, r.url, "")
	assert.Equal(t, r.compiledPath, ".")
	assert.Equal(t, len(r.params), 0)
}

func TestFullRequest(t *testing.T) {
	_, r := newRequest("/users", map[string]bool{"users": true})
	assert.Equal(t, r.compiledPath, "./users")
	assert.Equal(t, len(r.params), 0)
}

func TestRequestWithResourceAndParams(t *testing.T) {
	_, r := newRequest(
		"/users/123/sales/444",
		map[string]bool{
			"users": true,
			"sales": true,
		})

	assert.Equal(t, r.compiledPath, "./users/$/sales/$")
	assert.Equal(t, len(r.params), 2)
}

func createContext(method, URI string) *rocket.Context {
	req, _ := http.NewRequest(method, URI, nil)
	req.RequestURI = URI
	return rocket.NewContext(nil, req, 0, 0)
}

func TestRouterGet(t *testing.T) {
	handlers := map[string]int{}
	router := NewRouter()

	//handlers["/"] = false
	router.Add("GET", "/", func (ctx *rocket.Context)  {handlers["/"]++})
	router.Add("GET", "/users", func (ctx *rocket.Context)  {handlers["/users"]++})
	router.Add("GET", "/users/:userId", func (ctx *rocket.Context)  {handlers["/users/:userId"]++})
	router.Add("GET", "/users/:userId/sales", func (ctx *rocket.Context)  {handlers["/users/:userId/sales"]++})
	router.Add("GET", "/users/:userId/sales/:saleId", func (ctx *rocket.Context)  {handlers["/users/:userId/sales/:saleId"]++})

	router.handle(createContext("GET", ""))
	router.handle(createContext("GET", "/users"))
	router.handle(createContext("GET", "/users/1"))
	router.handle(createContext("GET", "/users/1/?order=asc"))
	router.handle(createContext("GET", "/users/sales"))
	router.handle(createContext("GET", "/users/_sales"))
	router.handle(createContext("GET", "/users/1/sales"))
	router.handle(createContext("GET", "/users/1/sales/2"))

	assert.Equal(t, 1, handlers["/"])
	assert.Equal(t, 1, handlers["/users"])
	assert.Equal(t, 3, handlers["/users/:userId"])
	assert.Equal(t, 1, handlers["/users/:userId/sales"])
	assert.Equal(t, 1, handlers["/users/:userId/sales/:saleId"])

}
