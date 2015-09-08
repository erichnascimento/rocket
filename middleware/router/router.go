package router

import (
	"errors"
	"log"
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

type Router struct {
	next      middleware.HandleFunc
	routes    map[string][]*Route
	resources map[string]bool
	root      string
}

func NewRouter(root string) *Router {
	return &Router{
		routes:    map[string][]*Route{},
		root:      root,
		resources: map[string]bool{},
	}
}

func (this *Router) CreateHandle(next middleware.HandleFunc) middleware.HandleFunc {
	this.next = next
	return this.handle
}

func (this *Router) handle(ctx *rocket.Context) {
	err, req := newRequest(this.root, ctx.Request.RequestURI, this.resources)
	if err == ErrorRequestHasDiferentRoot {
		this.next(ctx)
		return
	}

	for _, route := range this.routes[ctx.Request.Method] {
		log.Println(route.compiledRoute, req.compiledPath)
		if route.compiledRoute == req.compiledPath {
			if route.handler != nil {
				route.handler(&Context{ctx, route, req})
			}
			break
		}
	}

	this.next(ctx)
}

func (r *Router) Add(method, path string, handler HandleFunc) *Router {
	_, route := newRoute(path, handler)

	for k, v := range route.resources {
		r.resources[k] = v
	}

	// create an array for method
	if r.routes[method] == nil {
		r.routes[method] = make([]*Route, 0)
	}

	r.routes[method] = append(r.routes[method], route)

	return r
}

func getPath(path string) string {
	if i := strings.IndexAny(path, "?#"); i >= 0 {
		return path[:i]
	}

	return path
}

func explodePathParts(path string) []string {
	if path == "" {
		return []string{}
	}
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

var ErrorRequestHasDiferentRoot = errors.New("Request has diferent root")

type Request struct {
	url          string
	resources    map[string]bool
	compiledPath string
	params       []string
	root         string
}

func (r *Request) ParseURL(url string) error {
	r.compiledPath = "."
	r.url = url

	if r.root != "" {
		if !strings.HasPrefix(r.url, r.root) {
			return ErrorRequestHasDiferentRoot
		}

		// remove root from url
		r.url = strings.Replace(r.url, r.root, "", -1)
	}

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

func newRequest(root, url string, resources map[string]bool) (error, *Request) {
	r := new(Request)
	r.resources = resources
	r.params = make([]string, 0)
	r.root = root

	if err := r.ParseURL(url); err != nil {
		return err, nil
	}

	return nil, r
}
