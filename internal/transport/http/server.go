package http

import (
	"encoding/json"
	"net/http"
)

func Start() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	go http.ListenAndServe(":8080", nil)
}
