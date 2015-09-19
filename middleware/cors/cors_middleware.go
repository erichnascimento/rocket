// http://www.html5rocks.com/en/tutorials/cors/?redirect_from_locale=pt

package cors

import (
	"strconv"
	"strings"

	"github.com/erichnascimento/rocket"
	"github.com/erichnascimento/rocket/middleware"
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

// jsonBody is a middleware for handle cors
type Cors struct {
	config *Config
	next   middleware.HandleFunc
}

// CreateHandle create a new handler
func (c *Cors) CreateHandle(next middleware.HandleFunc) middleware.HandleFunc {
	c.next = next
	if c.config.Preflight {
		return c.PreflightHandle
	}

	return c.SimpleHandle
}

func (c *Cors) SimpleHandle(ctx *rocket.Context) {
	ctx.Header().Add("Access-Control-Allow-Credentials", strconv.FormatBool(c.config.AllowCredentials))
	ctx.Header().Add("Access-Control-Allow-Origin", c.config.AllowOrigin)

	if c.config.AllowMethods != nil {
		ctx.Header().Add("Access-Control-Allow-Methods", strings.Join(c.config.AllowMethods, ", "))
	}

	if c.config.AllowHeaders != nil {
		ctx.Header().Add("Access-Control-Allow-Headers", strings.Join(c.config.AllowHeaders, ", "))
	}
	//Access-Control-Allow-Methods
	c.next(ctx)
}

func (c *Cors) PreflightHandle(ctx *rocket.Context) {
	//ctx.Header().Add("Access-Control-Allow-Credentials", strconv.FormatBool(c.config.AllowCredentials))
	//ctx.Header().Add("Access-Control-Allow-Origin", c.config.AllowOrigin)
	//ctx.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	//Access-Control-Allow-Methods
	ctx.Header().Add("Access-Control-Allow-Methods", strings.Join(c.config.AllowMethods, ", "))
	ctx.Header().Add("Access-Control-Allow-Headers", strings.Join(c.config.AllowHeaders, ", "))
	c.next(ctx)
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
