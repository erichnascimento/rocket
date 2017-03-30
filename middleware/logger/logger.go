package logger

import (
	"log"
	"time"

	"github.com/dustin/go-humanize"
	"net/http"
	"github.com/erichnascimento/rocket/server"
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
	l.start = time.Now()
	l.next(rw, req)
	duration := time.Since(l.start) / time.Millisecond
	contentLength := humanize.Bytes(uint64(server.GetContentLength(rw)))
	log.Printf("%s %s %d %dms - %s", req.Method, req.RequestURI, server.GetStatusCode(rw), duration, contentLength)
}

// Create a new logger middleware
func NewLogger() *Logger {
	return &Logger{}
}
