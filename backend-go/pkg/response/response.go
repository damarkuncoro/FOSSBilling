package response

import (
	"encoding/json"
	"net/http"
)

type Meta struct {
	Page         int `json:"page,omitempty"`
	Limit        int `json:"limit,omitempty"`
	Offset       int `json:"offset,omitempty"`
	Total        int `json:"total,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
	TotalPages   int `json:"total_pages,omitempty"`
}


type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   *APIError   `json:"error"`
	Meta    *Meta       `json:"meta,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
		Error:   nil,
		Meta:    meta,
	})
}

func Error(w http.ResponseWriter, status int, code, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{
		Success: false,
		Data:    nil,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
