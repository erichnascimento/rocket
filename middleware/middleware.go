package middleware

import (
	"github.com/erichnascimento/rocket"
)

type HandleFunc func(ctx *rocket.Context)

type Middleware interface {
	CreateHandle(next HandleFunc) HandleFunc
}

type DynaHandleFunc func(ctx *rocket.Context, next HandleFunc)

type dynamiddleware struct {
	handler DynaHandleFunc
}

func NewDynamiddleware(handler DynaHandleFunc) Middleware {
	return &dynamiddleware{handler}
}

func (d *dynamiddleware) CreateHandle(next HandleFunc) HandleFunc {
	return func(ctx *rocket.Context) {
		d.handler(ctx, next)
	}
}
