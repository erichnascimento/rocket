package jsonbody

import (
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/erichnascimento/rocket"
	"github.com/erichnascimento/rocket/middleware"
)

func createJsonBody(next middleware.HandleFunc) *jsonBody {
	jb := NewJsonBody()

	if next == nil {
		next = func(ctx *rocket.Context) {}
	}

	jb.CreateHandle(next)

	return jb
}

func createContext(bodyContent, method, url string) *rocket.Context {
	reader := strings.NewReader(bodyContent)
	req, _ := http.NewRequest(method, url, reader)
	ctx := rocket.NewContext(nil, req, 0, 0)
	return ctx
}

func TestGetEmptyBody(t *testing.T) {
	jb := createJsonBody(nil)
	ctx := createContext(`""`, "GET", "")

	jb.handle(ctx)
	body, err := Get(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", body)
}

func TestGetSimpleStringBody(t *testing.T) {
	jb := createJsonBody(nil)
	ctx := createContext(`"this is a simple string"`, "GET", "")

	jb.handle(ctx)
	body, err := Get(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, "this is a simple string", body)
}

func TestGetNumericBody(t *testing.T) {
	jb := createJsonBody(nil)
	ctx := createContext(`123`, "GET", "")

	jb.handle(ctx)
	body, err := Get(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(123), body)
}

func TestGetObjectBody(t *testing.T) {
	jb := createJsonBody(nil)
	ctx := createContext(`{"name": "Ana"}`, "GET", "")

	jb.handle(ctx)
	body, err := Get(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Ana", body.(map[string]interface{})["name"])
}

func TestGetEmptyArrayBody(t *testing.T) {
	jb := createJsonBody(nil)
	ctx := createContext(`[]`, "GET", "")

	jb.handle(ctx)
	body, err := Get(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, []interface{}{}, body)
}

func TestGetFilledArrayBody(t *testing.T) {
	jb := createJsonBody(nil)
	ctx := createContext(`["", 123, "abc", {}]`, "GET", "")

	jb.handle(ctx)
	body, err := Get(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, []interface{}{"", float64(123), "abc", map[string]interface{}{}}, body)
}

func TestGetBooleanBody(t *testing.T) {
	jb := createJsonBody(nil)
	ctx := createContext(`true`, "GET", "")

	jb.handle(ctx)
	body, err := Get(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, body)
}

func TestDecodeEmptyObjectBody(t *testing.T) {
	jb := createJsonBody(nil)
	ctx := createContext(`{}`, "GET", "")

	jb.handle(ctx)
	var dest struct{}
	err := Decode(ctx, &dest)
	assert.Equal(t, nil, err)
	assert.Equal(t, struct{}{}, dest)
}

func TestDecodeNonEmptyObjectBody(t *testing.T) {
	jb := createJsonBody(nil)
	ctx := createContext(`{"id": "111", "name": "Jacob"}`, "GET", "")

	jb.handle(ctx)
	var dest struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	err := Decode(ctx, &dest)
	log.Printf("%v", dest)
	assert.Equal(t, nil, err)
	assert.Equal(t, "111", dest.Id)
	assert.Equal(t, "Jacob", dest.Name)
}
