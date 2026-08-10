package handlers

import (
	"net/http"
)

// HealthCheck godoc
// @Summary Health check
// @Description Check whether the DevSync API is running.
// @Tags Health
// @Produce json
// @Success 200 {object} utils.Response
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("API is running"))

}
