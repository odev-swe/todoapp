package libs

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status bool, code int, msg string, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	res := &Response{
		Message: msg,
		Status:  status,
		Code:    code,
		Data:    data,
	}

	return json.NewEncoder(w).Encode(res)
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	// handle time
	return json.NewDecoder(r.Body).Decode(v)
}
