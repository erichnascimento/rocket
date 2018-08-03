package response

import (
	"encoding/json"
	"net/http"
)

func SendJSON(rw http.ResponseWriter, v interface{}, statusCode int) error {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	return json.NewEncoder(rw).Encode(v)
}
