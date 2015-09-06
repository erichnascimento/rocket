package jsonbody

import (
	"encoding/json"

	"github.com/erichnascimento/rocket"
	"github.com/erichnascimento/rocket/middleware"
)

var (
	contextMap = map[*rocket.Context]interface{}{}
)

func Decode(ctx *rocket.Context, dest interface{}) error {
	return json.NewDecoder(ctx.Request.Body).Decode(&dest)
}

func Get(ctx *rocket.Context) (interface{}, error) {
	if val := contextMap[ctx]; val == nil {
		var body interface{}
		err := Decode(ctx, &body)
		if err != nil {
			return nil, err
		}
		if body == nil {
			body = ""
		}
		contextMap[ctx] = body
	}

	return contextMap[ctx], nil
}

// jsonBody is a middleware for parse json body content
type jsonBody struct {
	next middleware.HandleFunc
}

// CreateHandle create a new handler
func (j *jsonBody) CreateHandle(next middleware.HandleFunc) middleware.HandleFunc {
	j.next = next
	return j.handle
}

func (j *jsonBody) handle(ctx *rocket.Context) {
	//l.start = time.Now()
	contextMap[ctx] = nil
	j.next(ctx)
}

// NewJsonBody Create a new logger middleware
func NewJsonBody() *jsonBody {
	return &jsonBody{}
}
