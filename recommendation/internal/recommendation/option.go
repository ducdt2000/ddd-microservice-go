package recommendation

import "github.com/Rhymond/go-money"

type Option struct {
	HotelName     string
	Location      string
	PricePerNight money.Money
}
