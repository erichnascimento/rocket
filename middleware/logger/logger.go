package logger

import (
	"log"
	"time"

	"github.com/erichnascimento/rocket"
	"github.com/erichnascimento/rocket/middleware"

	"github.com/dustin/go-humanize"
)

// Logger is a middleware for loggin requests
type Logger struct {
	next  middleware.HandleFunc
	start time.Time
}

// CreateHandle create a new handler
func (l *Logger) CreateHandle(next middleware.HandleFunc) middleware.HandleFunc {
	l.next = next
	return l.handle
}

func (l *Logger) handle(ctx *rocket.Context) {
	l.start = time.Now()
	l.next(ctx)
	duration := time.Since(l.start) / time.Millisecond
	contentLenght := humanize.Bytes(uint64(ctx.GetContentLength()))
	log.Printf("%s %s %d %dms - %s", ctx.Request.Method, ctx.Request.RequestURI, ctx.GetStatusCode(), duration, contentLenght)
}

// Create a new logger middleware
func NewLogger() *Logger {
	return &Logger{}
}
