package server

import (
	"net/http"

	"github.com/erichnascimento/rocket/middleware"
	"github.com/erichnascimento/rocket"
)

func createFinalHandler() middleware.HandleFunc {
	return func (ctx *rocket.Context) {
		// nope
	}
}

type Server struct {
	server *http.Server
	handler middleware.HandleFunc
	middlewares []middleware.Middleware
}

func New(addr string) *Server {
	s := &Server{
		middlewares: make([]middleware.Middleware, 0),
	}

	s.server = &http.Server{
		Addr: addr,
		Handler: s,
	}

	return s
}

func (s *Server) Use(m middleware.Middleware) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) Serve() error {
	s.handler = createFinalHandler()
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.handler = s.middlewares[i].CreateHandle(s.handler)
	}

	return s.server.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := rocket.NewContext(w, r, 200, 0)
	s.handler(ctx)
}

func (s *Server) Stop() {
	/*for range time.Tick(time.Second) {
		log.Printf("Waiting %d", s.dispatcher.pending)
		if s.dispatcher.pending <= 0 {
			break
		}
	}
	s.dispatcher.w.Wait()*/
}

