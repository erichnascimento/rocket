package response

import (
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	written int
	status  int
}

// Override in order to capture bytes written
func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n

	return n, err
}

// Override in order to capture status
func (rw *responseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
	rw.status = code
}

func WrapResponseWriter(rw http.ResponseWriter) http.ResponseWriter {
	return &responseWriter{ResponseWriter: rw}
}

func GetContentLength(rw http.ResponseWriter) int {
	return rw.(*responseWriter).written
}

func GetStatusCode(rw http.ResponseWriter) int {
	return rw.(*responseWriter).status
}
