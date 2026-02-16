package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, v interface{}, status int) {
	js, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
