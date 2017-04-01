// http://www.html5rocks.com/en/tutorials/cors/?redirect_from_locale=pt

package middleware

import (
	"strconv"
	"strings"
	"net/http"
)

type Config struct {
	AllowCredentials bool
	AllowOrigin      string
	AllowMethods     []string
	AllowHeaders     []string
	Preflight        bool
}

func NewConfig() *Config {
	return &Config{
		AllowCredentials: true,
		AllowOrigin:      "*",
	}
}

func NewConfigPreflight(methods, headers []string) *Config {
	c := NewConfig()
	c.Preflight = true
	c.AllowMethods = methods
	c.AllowHeaders = headers

	return c
}

type Cors struct {
	config *Config
	next   http.HandlerFunc
}

func (c *Cors) Mount(next http.HandlerFunc) http.HandlerFunc {
	c.next = next
	if c.config.Preflight {
		return c.preflightHandler
	}

	return c.simpleHandler
}

func (c *Cors) simpleHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Access-Control-Allow-Credentials", strconv.FormatBool(c.config.AllowCredentials))
	rw.Header().Add("Access-Control-Allow-Origin", c.config.AllowOrigin)

	if c.config.AllowMethods != nil {
		rw.Header().Add("Access-Control-Allow-Methods", strings.Join(c.config.AllowMethods, ", "))
	}

	if c.config.AllowHeaders != nil {
		rw.Header().Add("Access-Control-Allow-Headers", strings.Join(c.config.AllowHeaders, ", "))
	}
	//Access-Control-Allow-Methods
	c.next(rw, req)
}

func (c *Cors) preflightHandler(rw http.ResponseWriter, req *http.Request) {
	//ctx.Header().Add("Access-Control-Allow-Credentials", strconv.FormatBool(c.config.AllowCredentials))
	//ctx.Header().Add("Access-Control-Allow-Origin", c.config.AllowOrigin)
	//ctx.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	//Access-Control-Allow-Methods
	rw.Header().Add("Access-Control-Allow-Methods", strings.Join(c.config.AllowMethods, ", "))
	rw.Header().Add("Access-Control-Allow-Headers", strings.Join(c.config.AllowHeaders, ", "))
	c.next(rw, req)
}

// NewJsonBody Create a new logger middleware
func NewCors(config *Config) *Cors {
	if config == nil {
		config = NewConfig()
	}

	return &Cors{
		config: config,
	}
}
