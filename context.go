package rocket

import (
	"net/http"
)

// Wrap the request and response
type Context struct {
	http.ResponseWriter
	Request *http.Request
	written int
	status  int
}

// Override to capture status
func (c *Context) WriteHeader(code int) {
	c.status = code
	c.ResponseWriter.WriteHeader(code)
}

// Override to capture written bytes
func (c *Context) Write(b []byte) (int, error) {
	n, err := c.ResponseWriter.Write(b)
	c.written += n

	return n, err
}

func (c *Context) GetContentLength() int {
	return c.written
}

func (c *Context) GetStatusCode() int {
	return c.status
}

func NewContext(w http.ResponseWriter, r *http.Request, status int, written int) *Context {
	return &Context{w, r, written, status}
}
