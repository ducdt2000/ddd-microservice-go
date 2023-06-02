package transport

import (
	"ddd-microservice-go/recommendation/internal/recommendation"
	"net/http"

	"github.com/gorilla/mux"
)

func NewMux(recommendationHandler recommendation.Handler) *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/recommendation", recommendationHandler.GetRecommendation).Methods(http.MethodGet)
	return m
}
