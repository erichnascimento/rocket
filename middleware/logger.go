package middleware

import (
	"errors"
	"log"
	"time"

	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/erichnascimento/rocket/server/response"
)

// Logger is a middleware for loggin requests
type Logger struct {
	next      http.HandlerFunc
	logFunc   LoggerFunc
}

// CreateHandle create a new handler
func (l *Logger) Mount(next http.HandlerFunc) http.HandlerFunc {
	l.next = next
	return l.handle
}

func (l *Logger) handle(rw http.ResponseWriter, req *http.Request) {
	rw = response.WrapResponseWriter(rw)
	start := time.Now()
	l.next(rw, req)
	duration := time.Since(start)// / time.Millisecond
	l.logFunc(rw, req, duration)
}

func defaultLogFunc(rw http.ResponseWriter, req *http.Request, duration time.Duration) {
	contentLength := humanize.Bytes(uint64(response.GetContentLength(rw)))
	log.Printf("%s %s %d %dms - %s", req.Method, req.RequestURI, response.GetStatusCode(rw), duration / time.Millisecond	, contentLength)
}

// Create a new logger middleware
func NewLogger() *Logger {
	return &Logger{
		logFunc: defaultLogFunc,
	}
}

// Create a new logger middleware
func NewCustomLogger(logFunc LoggerFunc) (*Logger, error) {
	if logFunc == nil {
		return nil, errors.New("logFunc can not be nil")
	}
	customLogger := NewLogger()
	customLogger.logFunc = logFunc
	return customLogger, nil
}

type LoggerFunc = func(rw http.ResponseWriter, req *http.Request, duration time.Duration)
