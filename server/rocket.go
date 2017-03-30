package server

import (
	"net/http"

	"github.com/erichnascimento/rocket/middleware"
)

type Rocket struct {
	server      *http.Server
	handler     http.HandlerFunc
	middlewares []middleware.Middleware
}

func (s *Rocket) ListenAndServe(addr string) error {
	s.handler = s.finalHandlerFunc
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.handler = s.middlewares[i].Mount(s.handler)
	}

	if s.server == nil {
		s.server = &http.Server{Handler: s}
	}
	s.server.Addr = addr

	return s.server.ListenAndServe()
}

func (s *Rocket) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	wrw := wrapResponseWriter(rw)
	s.handler(wrw, rq)
}

func (s *Rocket) Use(m middleware.Middleware) {
	if s.middlewares == nil {
		s.middlewares = make([]middleware.Middleware, 0)
	}
	s.middlewares = append(s.middlewares, m)
}

func (s *Rocket) finalHandlerFunc(rw http.ResponseWriter, rq *http.Request) {
	//log.Println("Done ", rq.Context())
}

func NewRocket() *Rocket {
	return new(Rocket)
}
