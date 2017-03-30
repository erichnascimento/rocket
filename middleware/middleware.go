package middleware

import (
	"github.com/erichnascimento/rocket"
	"net/http"
)

type HandleFunc func(ctx *rocket.Context)

type Middleware interface {
	Mount(next http.HandlerFunc) http.HandlerFunc
}