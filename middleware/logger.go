package middleware

import (
	"log"
	"time"

	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/erichnascimento/rocket/server/response"
)

// Logger is a middleware for loggin requests
type Logger struct {
	next  http.HandlerFunc
	start time.Time
}

// CreateHandle create a new handler
func (l *Logger) Mount(next http.HandlerFunc) http.HandlerFunc {
	l.next = next
	return l.handle
}

func (l *Logger) handle(rw http.ResponseWriter, req *http.Request) {
	rw = response.WrapResponseWriter(rw)

	l.start = time.Now()
	l.next(rw, req)
	duration := time.Since(l.start) / time.Millisecond
	contentLength := humanize.Bytes(uint64(response.GetContentLength(rw)))
	log.Printf("%s %s %d %dms - %s", req.Method, req.RequestURI, response.GetStatusCode(rw), duration, contentLength)
}

// Create a new logger middleware
func NewLogger() *Logger {
	return &Logger{}
}
