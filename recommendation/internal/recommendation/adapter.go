package recommendation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Rhymond/go-money"
)

type PartnershipAdaptor struct {
	client *http.Client
	url    string
}

func NewPartnershipAdaptor(client *http.Client, url string) (*PartnershipAdaptor, error) {
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}
	if url == "" {
		return nil, errors.New("url cannot be empty")
	}
	return &PartnershipAdaptor{client: client, url: url}, nil
}

func (pa PartnershipAdaptor) GetAvailability(
	ctx context.Context,
	tripStart time.Time,
	tripEnd time.Time,
	location string) ([]Option, error) {
	from := fmt.Sprintf("%d-%d-%d", tripStart.Year(), tripStart.Month(), tripStart.Day())
	to := fmt.Sprintf("%d-%d-%d", tripEnd.Year(), tripEnd.Month(), tripEnd.Day())
	url := fmt.Sprintf("%s/partnerships?location=%s&from=%s&to=%s", pa.url, location, from, to)
	res, err := pa.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call partnerships: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad request to partnerships: %d", res.StatusCode)
	}

	var partnershipRes PartnerShipsResponse
	if err := json.NewDecoder(res.Body).Decode(&partnershipRes); err != nil {
		return nil, fmt.Errorf("could not decode the response body of partnerships: %w", err)
	}

	options := make([]Option, len(partnershipRes.AvailableHotels))
	for i, v := range partnershipRes.AvailableHotels {
		options[i] = Option{
			HotelName:     v.Name,
			Location:      location,
			PricePerNight: *money.New(int64(v.PriceInUSDPerNight), "USD"),
		}
	}
	return options, nil
}
