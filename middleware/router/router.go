package router

import (
	"strings"

	"github.com/erichnascimento/rocket"
	"github.com/erichnascimento/rocket/middleware"
	//"log"
)

// Context type
type Context struct {
	*rocket.Context
	route   *Route
	request *Request
}

// GetParam return a param value
func (c *Context) GetParam(name string) string {
	for k, v := range c.route.params {
		if v == name {
			return c.request.params[k]
		}
	}
	return ""
}

type HandleFunc func(ctx *Context)

// TODO: Refactory to hash table
type routeEntries struct {
	get     []*Route
	post    []*Route
	put     []*Route
	del     []*Route
	patch   []*Route
	options []*Route
	head    []*Route
}

type Router struct {
	next      middleware.HandleFunc
	routes    *routeEntries
	resources map[string]bool
}

func NewRouter() *Router {
	return &Router{
		routes: &routeEntries{
			make([]*Route, 0),
			make([]*Route, 0),
			make([]*Route, 0),
			make([]*Route, 0),
			make([]*Route, 0),
			make([]*Route, 0),
			make([]*Route, 0),
		},
		resources: map[string]bool{},
	}
}

func (this *Router) CreateHandle(next middleware.HandleFunc) middleware.HandleFunc {
	this.next = next
	return this.handle
}

func (this *Router) handle(ctx *rocket.Context) {
	_, req := newRequest(ctx.Request.RequestURI, this.resources)

	var routes []*Route
	switch ctx.Request.Method {
	case "GET":
		routes = this.routes.get
	case "POST":
		routes = this.routes.post
	case "PUT":
		routes = this.routes.put
	case "DELETE":
		routes = this.routes.del
	case "PATCH":
		routes = this.routes.patch
	case "OPTIONS":
		routes = this.routes.options
	case "HEAD":
		routes = this.routes.head
	}

	for _, route := range routes {
		if route.compiledRoute == req.compiledPath {
			if route.handler != nil {
				route.handler(&Context{ctx, route, req})
			}
			break
		}
	}

	if this.next != nil {
		this.next(ctx)
	}
}

func (r *Router) Add(method, path string, handler HandleFunc) *Router {
	_, route := newRoute(path, handler)

	for k, v := range route.resources {
		r.resources[k] = v
	}

	switch method {
	case "GET":
		r.routes.get = append(r.routes.get, route)
	case "POST":
		r.routes.post = append(r.routes.post, route)
	case "PUT":
		r.routes.put = append(r.routes.put, route)
	case "DELETE":
		r.routes.del = append(r.routes.del, route)
	case "PATCH":
		r.routes.patch = append(r.routes.patch, route)
	case "OPTIONS":
		r.routes.options = append(r.routes.options, route)
	case "HEAD":
		r.routes.head = append(r.routes.head, route)
	}

	return r
}

func getPath(path string) string {
	if i := strings.IndexAny(path, "?#"); i >= 0 {
		return path[:i]
	}

	return path
}

func explodePathParts(path string) []string {
	return strings.Split(getPath(path[1:]), "/")
}

type Route struct {
	route         string
	compiledRoute string
	resources     map[string]bool
	params        []string
	handler       HandleFunc
}

func (r *Route) CompileRoute(route string) error {
	r.compiledRoute = "."

	if route == "" {
		return nil
	}

	for _, part := range explodePathParts(route) {
		switch {
		case part == "":
			continue
		case part[0] == ':':
			r.compiledRoute += "/$"
			r.params = append(r.params, part[1:])
		default:
			r.compiledRoute += "/" + part
			if _, ok := r.resources[part]; !ok {
				r.resources[part] = true
			}
		}
	}

	return nil
}

func newRoute(route string, handler HandleFunc) (error, *Route) {
	r := new(Route)
	r.resources = map[string]bool{}
	r.params = make([]string, 0)
	r.handler = handler
	err := r.CompileRoute(route)

	return err, r
}

type Request struct {
	url          string
	resources    map[string]bool
	compiledPath string
	params       []string
}

func (r *Request) ParseURL(url string) error {
	r.compiledPath = "."

	if r.url == url {
		return nil
	}

	r.url = url
	for _, v := range explodePathParts(r.url) {
		if v == "" {
			continue
		}

		if _, ok := r.resources[v]; ok {
			r.compiledPath += "/" + v
			continue
		}

		// param
		r.compiledPath += "/$"
		r.params = append(r.params, v)
	}

	return nil
}

func newRequest(url string, resources map[string]bool) (error, *Request) {
	r := new(Request)
	r.resources = resources
	r.params = make([]string, 0)
	err := r.ParseURL(url)

	return err, r
}
