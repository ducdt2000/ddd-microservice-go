package recommendation

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Rhymond/go-money"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) (*Handler, error) {
	if service == (Service{}) {
		return nil, errors.New("service cannot be empty")
	}
	return &Handler{Service: service}, nil
}

func (h Handler) GetRecommendation(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	log.Println(req.URL.Query())
	location := q.Get("location")
	if location == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	from := q.Get("from")
	if from == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	to := q.Get("to")
	if to == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	budget := q.Get("budget")
	if budget == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const expectedFormat = "2006-01-02"
	formattedStart, err := time.Parse(expectedFormat, from)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	formattedEnd, err := time.Parse(expectedFormat, to)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := strconv.ParseInt(budget, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	budgetMon := money.New(b, "USD")

	recommendation, err := h.Service.Get(req.Context(), formattedStart, formattedEnd, location, *budgetMon)
	log.Println(recommendation)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(GetRecommendationResponse{
		HotelName: recommendation.HotelName,
		TotalCost: struct {
			Cost     int64  "json:\"cost\""
			Currency string "json:\"currency\""
		}{
			Cost:     recommendation.TripPrice.Amount(),
			Currency: "USD",
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}
