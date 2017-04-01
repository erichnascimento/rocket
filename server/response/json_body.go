package response

import (
	"encoding/json"
	"net/http"
)

func EncodeJSON(rw http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(rw).Encode(v)
}
