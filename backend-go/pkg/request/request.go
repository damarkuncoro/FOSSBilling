package request

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Decode maps JSON request body to a struct
func Decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// GetID extracts "id" from path parameters as int64
func GetID(r *http.Request) int64 {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id
}

// GetParam extracts a string parameter from path
func GetParam(r *http.Request, key string) string {
	return r.PathValue(key)
}

// GetLimitOffset extracts pagination parameters from query string
func GetLimitOffset(r *http.Request) (int, int) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 { limit = 50 }
	if offset < 0 { offset = 0 }
	return limit, offset
}
