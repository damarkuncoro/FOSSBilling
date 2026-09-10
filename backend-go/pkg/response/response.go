package response

import (
	"encoding/json"
	"net/http"
)

type Meta struct { Total, Limit, Offset, Page, TotalPages int }
type Err struct { Code, Message string; Details any }
type Res struct { Success bool `json:"success"`; Data any `json:"data"`; Error *Err `json:"error"`; Meta *Meta `json:"meta,omitempty"` }

func JSON(w http.ResponseWriter, s int, d any, m *Meta) {
	w.Header().Set("Content-Type", "application/json"); w.WriteHeader(s)
	_ = json.NewEncoder(w).Encode(Res{Success: true, Data: d, Meta: m})
}

func Error(w http.ResponseWriter, s int, c, m string, d any) {
	w.Header().Set("Content-Type", "application/json"); w.WriteHeader(s)
	_ = json.NewEncoder(w).Encode(Res{Success: false, Error: &Err{c, m, d}})
}
