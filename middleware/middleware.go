package middleware

import (
	"net/http"
)


type Middleware interface {
	Mount(next http.HandlerFunc) http.HandlerFunc
}