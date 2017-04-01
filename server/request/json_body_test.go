package request

import (
	"net/http"
	"testing"
	"strings"
)

func TestDecodeJSON(t *testing.T) {
	body := strings.NewReader(`{"foo": "bar"}`)
	req, _ := http.NewRequest(http.MethodPost, "foo_url", body)
	var v map[string]interface{}
	err := DecodeJSON(req, &v)
	if err != nil {
		t.Error(err)
	}

	if v["foo"] != "bar" {
		t.Errorf(`Error decoding json body: expected "foo=bar", given "foo=%s"`, v["foo"])
	}
}
