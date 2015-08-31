package middleware

import (
	"log"
	"time"

	"github.com/erichnascimento/rocket"

	"github.com/dustin/go-humanize"
)

type Logger struct {
	next HandleFunc
	start time.Time
}

func (this *Logger) CreateHandle(next HandleFunc) HandleFunc {
	this.next = next
	return this.handle
}

func (this *Logger) handle(ctx *rocket.Context) {
	this.start = time.Now()
	this.next(ctx)
	duration := time.Since(this.start) / time.Millisecond
	contentLenght := humanize.Bytes(uint64(ctx.GetContentLength()))
	log.Printf("%s %s %d %dms - %s", ctx.Request.Method, ctx.Request.RequestURI, ctx.GetStatusCode(), duration, contentLenght)
}
