package server

import (
	"net/http"

	"github.com/erichnascimento/rocket/middleware"
)

type Server struct {
	server      *http.Server
	handler     http.HandlerFunc
	middlewares []middleware.Middleware
}

func (s *Server) ListenAndServe(addr string) error {
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

func (s *Server) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	s.handler(rw, rq)
}

func (s *Server) Use(m middleware.Middleware) {
	if s.middlewares == nil {
		s.middlewares = make([]middleware.Middleware, 0)
	}
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) finalHandlerFunc(rw http.ResponseWriter, rq *http.Request) {
	//log.Println("Done ", rq.Context())
}

func NewServer() *Server {
	return new(Server)
}
