package middleware

import (
	"errors"
	"strings"
	"net/http"
	"context"
)


type Router struct {
	next      http.HandlerFunc
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

func (r *Router) Mount(next http.HandlerFunc) http.HandlerFunc {
	r.next = next
	return r.handle
}

func (r *Router) handle(rw http.ResponseWriter, req *http.Request) {
	err, reqInfo := newRequestInfo(r.root, req.RequestURI, r.resources)
	if err == ErrorRequestHasDifferentRoot {
		r.next(rw, req)
		return
	}

	for _, route := range r.routes[req.Method] {
		if route.compiledRoute == reqInfo.compiledPath {
			ctx := req.Context()
			for k, name := range route.params {
				ctx = context.WithValue(ctx, name, reqInfo.params[k])
			}

			reqWithParams := req.WithContext(ctx)
			for _, handler := range  route.handlers {
				handler(rw, reqWithParams)
			}
			break
		}
	}

	r.next(rw, req)
}

func (r *Router) Add(method, path string, handlers...http.HandlerFunc) *Router {
	_, route := newRoute(path, handlers...)

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

func (r *Router) Get(path string, handlers...http.HandlerFunc) *Router {
	return r.Add(http.MethodGet, path, handlers...)
}

func (r *Router) Post(path string, handlers...http.HandlerFunc) *Router {
	return r.Add(http.MethodPost, path, handlers...)
}

func (r *Router) Put(path string, handlers...http.HandlerFunc) *Router {
	return r.Add(http.MethodPut, path, handlers...)
}

func (r *Router) Patch(path string, handlers...http.HandlerFunc) *Router {
	return r.Add(http.MethodPatch, path, handlers...)
}

func (r *Router) Delete(path string, handlers...http.HandlerFunc) *Router {
	return r.Add(http.MethodDelete, path, handlers...)
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
	handlers	[]http.HandlerFunc
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

func newRoute(route string, handlers...http.HandlerFunc) (error, *Route) {
	r := new(Route)
	r.resources = map[string]bool{}
	r.params = make([]string, 0)
	r.handlers = handlers
	err := r.CompileRoute(route)

	return err, r
}

var ErrorRequestHasDifferentRoot = errors.New("Request has diferent root")

type RequestInfo struct {
	url          string
	resources    map[string]bool
	compiledPath string
	params       []string
	root         string
}

func (r *RequestInfo) ParseURL(url string) error {
	r.compiledPath = "."
	r.url = url

	if r.root != "" {
		if !strings.HasPrefix(r.url, r.root) {
			return ErrorRequestHasDifferentRoot
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

func newRequestInfo(root, url string, resources map[string]bool) (error, *RequestInfo) {
	r := new(RequestInfo)
	r.resources = resources
	r.params = make([]string, 0)
	r.root = root
	if err := r.ParseURL(url); err != nil {
		return err, nil
	}

	return nil, r
}