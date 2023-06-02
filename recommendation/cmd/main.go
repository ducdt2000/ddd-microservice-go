package main

import (
	"ddd-microservice-go/recommendation/internal/recommendation"
	"ddd-microservice-go/recommendation/internal/transport"
	"log"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

func main() {
	c := retryablehttp.NewClient()
	c.RetryMax = 10

	partnerAdaptor, err := recommendation.NewPartnershipAdaptor(
		c.StandardClient(),
		"http://localhost:3031",
	)
	if err != nil {
		log.Fatal("failed to create a partnerAdaptor: ", err)
	}

	log.Println(*partnerAdaptor)

	service, err := recommendation.NewService(partnerAdaptor)
	if err != nil {
		log.Fatal("failed to create a service: ", err)
	}

	handler, err := recommendation.NewHandler(*service)
	if err != nil {
		log.Fatal("failed to create a handler: ", err)
	}

	m := transport.NewMux(*handler)

	if err := http.ListenAndServe(":4040", m); err != nil {
		log.Fatal("server errored: ", err)
	}
}
