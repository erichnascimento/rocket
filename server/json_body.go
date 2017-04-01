package server

import (
	"encoding/json"
	"net/http"
)

func Decode(req *http.Request, v interface{}) error {
	defer req.Body.Close()

	return json.NewDecoder(req.Body).Decode(&v)
}
