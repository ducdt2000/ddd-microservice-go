package recommendation

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Rhymond/go-money"
)

type Service struct {
	AvailabilityGetter AvailabilityGetter
}

func NewService(availabilityGetter AvailabilityGetter) (*Service, error) {
	if availabilityGetter == nil {
		return nil, errors.New("availability must not be nil")
	}
	return &Service{AvailabilityGetter: availabilityGetter}, nil
}

func (s Service) Get(
	ctx context.Context,
	tripStart time.Time,
	tripEnd time.Time,
	location string,
	budget money.Money,
) (*Recommendation, error) {
	switch {
	case tripStart.IsZero():
		return nil, errors.New("trip start cannot be empty")
	case tripEnd.IsZero():
		return nil, errors.New("trip end cannot be empty")
	case location == "":
		return nil, errors.New("location cannot be empty")
	}

	options, err := s.AvailabilityGetter.GetAvailability(ctx, tripStart, tripEnd, location)
	if err != nil {
		return nil, fmt.Errorf("error getting availability: %w", err)
	}

	tripDuration := math.Round(tripEnd.Sub(tripStart).Hours() / 24)
	lowestPrice := money.NewFromFloat(math.MaxInt, "USD")
	var cheapestTrip *Option
	for _, option := range options {
		price := option.PricePerNight.Multiply(int64(tripDuration))
		if ok, _ := price.GreaterThanOrEqual(&budget); ok {
			continue
		}
		if ok, _ := price.LessThan(&budget); ok {
			lowestPrice = price
			cheapestTrip = &option
		}
	}
	if cheapestTrip == nil {
		return nil, errors.New("no trips within budget")
	}
	return &Recommendation{
		TripStart: tripStart,
		TripEnd:   tripEnd,
		HotelName: cheapestTrip.HotelName,
		Location:  cheapestTrip.Location,
		TripPrice: *lowestPrice,
	}, nil
}
