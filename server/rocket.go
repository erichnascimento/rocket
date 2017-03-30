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

func (s *Rocket) finalHandlerFunc(rw http.ResponseWriter, rq *http.Request) {
	//log.Println("Done ", rq.Context())
}

func (s *Rocket) Use(m middleware.Middleware) {
	if s.middlewares == nil {
		s.middlewares = make([]middleware.Middleware, 0)
	}
	s.middlewares = append(s.middlewares, m)
}

type ResponseWriter struct {
	http.ResponseWriter
	written int
	status int
}

// Override in order to capture bytes written
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n

	return n, err
}

// Override in order to capture status
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
	rw.status = code
}

func wrapResponseWriter(rw http.ResponseWriter) http.ResponseWriter {
	result := &ResponseWriter{ ResponseWriter: rw}
	result.status = http.StatusOK
	rw.WriteHeader(result.status)

	return result
}

func GetContentLength(rw http.ResponseWriter) int {
	return rw.(*ResponseWriter).written
}

func GetStatusCode(rw http.ResponseWriter) int {
	return rw.(*ResponseWriter).status
}