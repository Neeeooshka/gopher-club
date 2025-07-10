package httputil

import (
	"encoding/json"
	"log"
	"net/http"
)

// WriteJSON write required headers, status and HTTP response body in JSON format
func WriteJSON(w http.ResponseWriter, v interface{}) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	jsonData, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("failed to write body: %v", err)
	}
}
