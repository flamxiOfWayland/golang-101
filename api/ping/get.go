package ping

import (
	"fmt"
	"net/http"
)

// @Produce  json
// @Success 200
// @Router /ping [get]
func GetHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get handler, status: %d\n", http.StatusOK)
}
