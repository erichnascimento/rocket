package server

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	written int
	status int
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

func wrapResponseWriter(rw http.ResponseWriter) http.ResponseWriter {
	result := &responseWriter{ ResponseWriter: rw}
	result.status = http.StatusOK
	rw.WriteHeader(result.status)

	return result
}

func GetContentLength(rw http.ResponseWriter) int {
	return rw.(*responseWriter).written
}

func GetStatusCode(rw http.ResponseWriter) int {
	return rw.(*responseWriter).status
}