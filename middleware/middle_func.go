package middleware

import "net/http"

type middleFunc struct {
	fn func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc)
}
func (mf middleFunc) Mount(next http.HandlerFunc) http.HandlerFunc {
	return func (rw http.ResponseWriter, req *http.Request) {
		mf.fn(rw, req, next)
	}
}

func NewMiddleFunc(fn func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc)) Middleware {
	return &middleFunc{fn}
}