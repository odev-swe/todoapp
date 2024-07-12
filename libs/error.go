package libs

import "net/http"

func InternalServerError(w http.ResponseWriter, msg string) {
	WriteJSON(w, false, http.StatusInternalServerError, msg, nil)
}

func NotFound(w http.ResponseWriter, msg string) {
	WriteJSON(w, false, http.StatusNotFound, msg, nil)
}

func BadRequest(w http.ResponseWriter, msg string) {
	WriteJSON(w, false, http.StatusBadRequest, msg, nil)
}

func Unauthorized(w http.ResponseWriter, msg string) {
	WriteJSON(w, false, http.StatusUnauthorized, msg, nil)
}
