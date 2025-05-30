package ping

import (
	"log/slog"
	"net/http"
)

// @Summary Ping
// @Description Ping Route
func Handler(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		return
	}
	if r.Method == http.MethodGet {
		GetHandle(w, r)
		return
	}

	if r.Method == http.MethodPost {
		PostHandle(w, r)
		return
	}

	slog.Warn("Ping Handler", "unsupported method", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}
