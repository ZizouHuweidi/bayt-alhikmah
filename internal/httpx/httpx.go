package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ReadJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message})
}

func Pagination(r *http.Request) (limit, offset int) {
	limit = 100
	offset = 0

	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > 100 {
		limit = 100
	}

	if value := r.URL.Query().Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}

func BearerToken(r *http.Request) (string, error) {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
		return "", errors.New("missing bearer token")
	}
	return header[len(prefix):], nil
}
