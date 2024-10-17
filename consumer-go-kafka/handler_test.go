package main

import (
	"testing"

	"github.com/pact-foundation/pact-go/v2/log"
	matchers "github.com/pact-foundation/pact-go/v2/matchers"
	message "github.com/pact-foundation/pact-go/v2/message/v4"
	"github.com/stretchr/testify/assert"
)

func TestConsumerV4(t *testing.T) {
	log.SetLogLevel("INFO")

	// 1 We write a small adapter that will take the incoming Message
	// and call the function with the correct type
	var handlerWrapper = func(m message.AsynchronousMessage) error {
		return handler(m.Body.(*ProductEvent).Product, m.Body.(*ProductEvent).Event)
	}

	// 2 Create the Pact Message Consumer
	mockProvider, err := message.NewAsynchronousPact(message.Config{
		Consumer: "PactGoV4KafkaConsumer",
		Provider: "pactflow-example-provider-js-kafka",
	})
	assert.NoError(t, err)

	// 3 Write the consumer test, and call ConsumedBy
	// passing through the function, and then Verify
	// to ensure it is correctly executed
	err = mockProvider.
		AddAsynchronousMessage().
		ExpectsToReceive("A foo").
		WithMetadata(map[string]string{"kafka_topic": "products"}).
		WithJSONContent(matchers.Map{
			"id":      matchers.Like("123"),
			"type":    matchers.Like("Product Range"),
			"name":    matchers.Like("Some Product"),
			"version": matchers.Like("v1"),
			"event":   matchers.Regex("UPDATED", "^(CREATED|UPDATED|DELETED)$"),
		}).AsType(&ProductEvent{}).
		ConsumedBy(handlerWrapper).Verify(t)
	assert.NoError(t, err)
}
