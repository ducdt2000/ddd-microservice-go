package recommendation

type GetRecommendationResponse struct {
	HotelName string `json:"hotelName"`
	TotalCost struct {
		Cost     int64  `json:"cost"`
		Currency string `json:"currency"`
	} `json:"totalCost"`
}
