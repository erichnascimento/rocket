package middleware

import (
	"github.com/erichnascimento/rocket"
	"net/http"
)


type Middleware interface {
	Mount(next http.HandlerFunc) http.HandlerFunc
}