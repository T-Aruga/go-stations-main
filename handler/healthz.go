package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		res := &model.HealthzResponse{Message: "OK"}
		encoder := json.NewEncoder(w)

		err := encoder.Encode(res)
		if err != nil {
			log.Println(err)
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
