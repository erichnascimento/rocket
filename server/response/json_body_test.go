package response

import (
	"testing"
	"net/http/httptest"
	"encoding/json"
)

func TestEncodeJSON(t *testing.T) {
	rw := httptest.NewRecorder()
	v := map[string]interface{}{
		"foo": "bar",
	}
	err := EncodeJSON(rw, v)
	if err != nil {
		t.Error(err)
	}

	var o map[string]interface{}
	err = json.Unmarshal(rw.Body.Bytes(), &o)
	if err != nil {
		t.Error(err)
	}
	if o["foo"] != `bar` {
		t.Errorf(`Error enconding json body: expected 'foo=bar', given "foo=%s"`, o["foo"])
	}
}
