package router

import (
	"net/http"
	"testing"

	"github.com/erichnascimento/rocket"

	"github.com/bmizerany/assert"
)

const DefaultRouteRoot = "/api/v1"

func TestGetPath(t *testing.T) {
	paths := map[string]string{
		"/user":          "/user",
		"/user?":         "/user",
		"/user?p=123":    "/user",
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

func TestRootPrefix(t *testing.T) {
	err, r := newRequest("/api/v1", "", nil)
	assert.Equal(t, err, ErrorRequestHasDiferentRoot)

	err, r = newRequest("/api/v1", "/bla/bla/bla", nil)
	assert.Equal(t, err, ErrorRequestHasDiferentRoot)

	err, r = newRequest("/api/v1", "/api/v1/user", map[string]bool{"user": true})
	assert.Equal(t, err, nil)
	assert.Equal(t, r.compiledPath, "./user")
}

func TestEmptyRequest(t *testing.T) {
	err, r := newRequest("", "", nil)

	assert.Equal(t, err, nil)
	assert.NotEqual(t, r, nil)
	assert.Equal(t, r.url, "")
	assert.Equal(t, r.compiledPath, ".")
	assert.Equal(t, len(r.params), 0)
}

func TestFullRequest(t *testing.T) {
	_, r := newRequest("/api/v1", "/api/v1/users", map[string]bool{"users": true})
	assert.Equal(t, r.compiledPath, "./users")
	assert.Equal(t, len(r.params), 0)
}

func TestRequestWithResourceAndParams(t *testing.T) {
	_, r := newRequest(
		"/api/v1",
		"/api/v1/users/123/sales/444",
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

func nopeHandleFunc(ctx *rocket.Context) {

}

func TestRouterGet(t *testing.T) {
	handlers := map[string]int{}
	router := NewRouter("")
	router.CreateHandle(nopeHandleFunc)

	//handlers["/"] = false
	router.Add("GET", "/", func(ctx *Context) { handlers["/"]++ })
	router.Add("GET", "/users", func(ctx *Context) { handlers["/users"]++ })
	router.Add("GET", "/users/:userId", func(ctx *Context) { handlers["/users/:userId"]++ })
	router.Add("GET", "/users/:userId/sales", func(ctx *Context) { handlers["/users/:userId/sales"]++ })
	router.Add("GET", "/users/:userId/sales/:saleId", func(ctx *Context) { handlers["/users/:userId/sales/:saleId"]++ })

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

func TestContext(t *testing.T) {
	router := NewRouter(DefaultRouteRoot)
	router.CreateHandle(nopeHandleFunc)

	// Unexistent param should return empty string
	router.Add("GET", "/", func(ctx *Context) {
		assert.Equal(t, "", ctx.GetParam("userId"))
	})

	router.Add("GET", "/users/:userId", func(ctx *Context) {
		assert.Equal(t, "1", ctx.GetParam("userId"))
	})

	router.Add("GET", "/users/:userId/sales/:saleId", func(ctx *Context) {
		assert.Equal(t, "1", ctx.GetParam("userId"))
		assert.Equal(t, "2", ctx.GetParam("saleId"))
	})

	router.handle(createContext("GET", "/"))
	router.handle(createContext("GET", "/users/1"))
	router.handle(createContext("GET", "/users/1/sales/2"))
}

func TestRouterPost(t *testing.T) {
	handlers := map[string]int{}
	router := NewRouter("")
	router.CreateHandle(nopeHandleFunc)

	//handlers["/"] = false
	router.Add("POST", "/", func(ctx *Context) { handlers["/"]++ })
	router.Add("POST", "/users", func(ctx *Context) { handlers["/users"]++ })
	router.Add("POST", "/users/:userId", func(ctx *Context) { handlers["/users/:userId"]++ })
	router.Add("POST", "/users/:userId/sales", func(ctx *Context) { handlers["/users/:userId/sales"]++ })
	router.Add("POST", "/users/:userId/sales/:saleId", func(ctx *Context) { handlers["/users/:userId/sales/:saleId"]++ })

	router.handle(createContext("POST", "/"))
	router.handle(createContext("POST", "/users"))
	router.handle(createContext("POST", "/users/1"))
	router.handle(createContext("POST", "/users/1/sales"))
	router.handle(createContext("POST", "/users/1/sales/123"))

	assert.Equal(t, 1, handlers["/"])
	assert.Equal(t, 1, handlers["/users"])
	assert.Equal(t, 1, handlers["/users/:userId"])
	assert.Equal(t, 1, handlers["/users/:userId/sales"])
	assert.Equal(t, 1, handlers["/users/:userId/sales/:saleId"])
}

func TestRouterPut(t *testing.T) {
	handlers := map[string]int{}
	router := NewRouter("")
	router.CreateHandle(nopeHandleFunc)

	//handlers["/"] = false
	router.Add("PUT", "/", func(ctx *Context) { handlers["/"]++ })
	router.Add("PUT", "/users", func(ctx *Context) { handlers["/users"]++ })
	router.Add("PUT", "/users/:userId", func(ctx *Context) { handlers["/users/:userId"]++ })
	router.Add("PUT", "/users/:userId/sales", func(ctx *Context) { handlers["/users/:userId/sales"]++ })
	router.Add("PUT", "/users/:userId/sales/:saleId", func(ctx *Context) { handlers["/users/:userId/sales/:saleId"]++ })

	router.handle(createContext("PUT", "/"))
	router.handle(createContext("PUT", "/users"))
	router.handle(createContext("PUT", "/users/1"))
	router.handle(createContext("PUT", "/users/1/sales"))
	router.handle(createContext("PUT", "/users/1/sales/123"))

	assert.Equal(t, 1, handlers["/"])
	assert.Equal(t, 1, handlers["/users"])
	assert.Equal(t, 1, handlers["/users/:userId"])
	assert.Equal(t, 1, handlers["/users/:userId/sales"])
	assert.Equal(t, 1, handlers["/users/:userId/sales/:saleId"])
}

func TestRouterDelete(t *testing.T) {
	handlers := map[string]int{}
	router := NewRouter("")
	router.CreateHandle(nopeHandleFunc)

	//handlers["/"] = false
	router.Add("DELETE", "/", func(ctx *Context) { handlers["/"]++ })
	router.Add("DELETE", "/users", func(ctx *Context) { handlers["/users"]++ })
	router.Add("DELETE", "/users/:userId", func(ctx *Context) { handlers["/users/:userId"]++ })
	router.Add("DELETE", "/users/:userId/sales", func(ctx *Context) { handlers["/users/:userId/sales"]++ })
	router.Add("DELETE", "/users/:userId/sales/:saleId", func(ctx *Context) { handlers["/users/:userId/sales/:saleId"]++ })

	router.handle(createContext("DELETE", "/"))
	router.handle(createContext("DELETE", "/users"))
	router.handle(createContext("DELETE", "/users/1"))
	router.handle(createContext("DELETE", "/users/1/sales"))
	router.handle(createContext("DELETE", "/users/1/sales/123"))

	assert.Equal(t, 1, handlers["/"])
	assert.Equal(t, 1, handlers["/users"])
	assert.Equal(t, 1, handlers["/users/:userId"])
	assert.Equal(t, 1, handlers["/users/:userId/sales"])
	assert.Equal(t, 1, handlers["/users/:userId/sales/:saleId"])
}

func TestRouterPatch(t *testing.T) {
	handlers := map[string]int{}
	router := NewRouter("")
	router.CreateHandle(nopeHandleFunc)

	//handlers["/"] = false
	router.Add("PATCH", "/", func(ctx *Context) { handlers["/"]++ })
	router.Add("PATCH", "/users", func(ctx *Context) { handlers["/users"]++ })
	router.Add("PATCH", "/users/:userId", func(ctx *Context) { handlers["/users/:userId"]++ })
	router.Add("PATCH", "/users/:userId/sales", func(ctx *Context) { handlers["/users/:userId/sales"]++ })
	router.Add("PATCH", "/users/:userId/sales/:saleId", func(ctx *Context) { handlers["/users/:userId/sales/:saleId"]++ })

	router.handle(createContext("PATCH", "/"))
	router.handle(createContext("PATCH", "/users"))
	router.handle(createContext("PATCH", "/users/1"))
	router.handle(createContext("PATCH", "/users/1/sales"))
	router.handle(createContext("PATCH", "/users/1/sales/123"))

	assert.Equal(t, 1, handlers["/"])
	assert.Equal(t, 1, handlers["/users"])
	assert.Equal(t, 1, handlers["/users/:userId"])
	assert.Equal(t, 1, handlers["/users/:userId/sales"])
	assert.Equal(t, 1, handlers["/users/:userId/sales/:saleId"])
}

func TestRouterOptions(t *testing.T) {
	handlers := map[string]int{}
	router := NewRouter("")
	router.CreateHandle(nopeHandleFunc)

	//handlers["/"] = false
	router.Add("OPTIONS", "/", func(ctx *Context) { handlers["/"]++ })
	router.Add("OPTIONS", "/users", func(ctx *Context) { handlers["/users"]++ })
	router.Add("OPTIONS", "/users/:userId", func(ctx *Context) { handlers["/users/:userId"]++ })
	router.Add("OPTIONS", "/users/:userId/sales", func(ctx *Context) { handlers["/users/:userId/sales"]++ })
	router.Add("OPTIONS", "/users/:userId/sales/:saleId", func(ctx *Context) { handlers["/users/:userId/sales/:saleId"]++ })

	router.handle(createContext("OPTIONS", "/"))
	router.handle(createContext("OPTIONS", "/users"))
	router.handle(createContext("OPTIONS", "/users/1"))
	router.handle(createContext("OPTIONS", "/users/1/sales"))
	router.handle(createContext("OPTIONS", "/users/1/sales/123"))

	assert.Equal(t, 1, handlers["/"])
	assert.Equal(t, 1, handlers["/users"])
	assert.Equal(t, 1, handlers["/users/:userId"])
	assert.Equal(t, 1, handlers["/users/:userId/sales"])
	assert.Equal(t, 1, handlers["/users/:userId/sales/:saleId"])
}

func TestRouterHead(t *testing.T) {
	handlers := map[string]int{}
	router := NewRouter("")
	router.CreateHandle(nopeHandleFunc)

	//handlers["/"] = false
	router.Add("HEAD", "/", func(ctx *Context) { handlers["/"]++ })
	router.Add("HEAD", "/users", func(ctx *Context) { handlers["/users"]++ })
	router.Add("HEAD", "/users/:userId", func(ctx *Context) { handlers["/users/:userId"]++ })
	router.Add("HEAD", "/users/:userId/sales", func(ctx *Context) { handlers["/users/:userId/sales"]++ })
	router.Add("HEAD", "/users/:userId/sales/:saleId", func(ctx *Context) { handlers["/users/:userId/sales/:saleId"]++ })

	router.handle(createContext("HEAD", "/"))
	router.handle(createContext("HEAD", "/users"))
	router.handle(createContext("HEAD", "/users/1"))
	router.handle(createContext("HEAD", "/users/1/sales"))
	router.handle(createContext("HEAD", "/users/1/sales/123"))

	assert.Equal(t, 1, handlers["/"])
	assert.Equal(t, 1, handlers["/users"])
	assert.Equal(t, 1, handlers["/users/:userId"])
	assert.Equal(t, 1, handlers["/users/:userId/sales"])
	assert.Equal(t, 1, handlers["/users/:userId/sales/:saleId"])
}
