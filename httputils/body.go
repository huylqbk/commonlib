package httputils

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response[T any] struct {
	StatusCode int   `json:"-"`
	Data       *T    `json:"data,omitempty"`
	Error      error `json:"error,omitempty"`
}

// ParseRequestBody parse request body using json new decoder instead on marshall/unmarshall
func ParseRequestBody[T comparable](req *http.Request) T {
	reqBody := json.NewDecoder(req.Body)
	var data T
	err := reqBody.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func Json(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
