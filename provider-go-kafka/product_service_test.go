package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pact-foundation/pact-go/v2/log"
	"github.com/pact-foundation/pact-go/v2/message"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
)

var pactDir, _ = os.Getwd()

func TestMessageProvider(t *testing.T) {
	log.SetLogLevel("INFO")

	verifier := provider.NewVerifier()

	// Map test descriptions to message producer (handlers)
	functionMappings := message.Handlers{
		"a product event update": func([]models.ProviderState) (message.Body, message.Metadata, error) {
			return createEvent(Product{
					ID:      "567",
					Name:    "A Product",
					Version: "v7",
					Type:    "Household Products",
				}, "UPDATED"), message.Metadata{
					"contentType": "application/json",
					"kafka_topic": "products",
				}, nil
		},
	}

	// Verify the Provider with local Pact Files

	verifier.VerifyProvider(t, provider.VerifyRequest{
		PactFiles:       []string{filepath.ToSlash(fmt.Sprintf("%s/../consumer-go-kafka/pacts/pactflow-example-consumer-go-kafka-pactflow-example-provider-go-kafka.json", pactDir))},
		Provider:        "pactflow-example-provider-go-kafka",
		MessageHandlers: functionMappings,
	})

}
