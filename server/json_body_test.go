package server

import (
	"net/http"
	"testing"
	"strings"
)

func TestDecode(t *testing.T) {
	body := strings.NewReader(`{"foo": "bar"}`)
	req, _ := http.NewRequest(http.MethodPost, "foo_url", body)
	var v map[string]interface{}
	err := Decode(req, &v)
	if err != nil {
		t.Error(err)
	}

	if v["foo"] != "bar" {
		t.Errorf(`Error decoding json body: expected "foo=bar", given "foo=%s"`, v["foo"])
	}
}
