package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendJSON(t *testing.T) {
	rw := httptest.NewRecorder()
	v := map[string]interface{}{
		"foo": "bar",
	}
	err := SendJSON(rw, v, http.StatusOK)
	if err != nil {
		t.Error(err)
	}

	if expected, given := http.StatusOK, rw.Result().StatusCode; expected != given {
		t.Errorf(`Invalid status code. Expected %d, given %d`, expected, given)
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
