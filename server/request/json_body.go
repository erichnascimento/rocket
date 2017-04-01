package request

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(req *http.Request, v interface{}) error {
	defer req.Body.Close()

	return json.NewDecoder(req.Body).Decode(&v)
}
