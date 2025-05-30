package ping

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

// @Produce  json
// @Param input body Default true "Default Input"
// @Success 200
// @Router /ping [post]
func PostHandle(w http.ResponseWriter, r *http.Request) {
	rawData, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("couldn't read data"))
		slog.Error("default, invalid data received", "err:", err)
		return
	}

	var data Default
	if err := json.Unmarshal(rawData, &data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid data format"))
		slog.Error("default, invalid data format", "err:", err)
		return
	}
	slog.Info("default request", "data:", data)
	return
}
