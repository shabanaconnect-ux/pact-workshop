const { MessageProviderPact, providerWithMetadata } = require("@pact-foundation/pact");
const { Product } = require("./product");
const { createEvent } = require("./product.event");
const path = require("path");

describe("Message provider tests", () => {

  const pactUrl = process.env.PACT_URL || path.join(__dirname, "..", "..", "..", "consumer-js-kafka", "pacts", "pactflow-example-consumer-js-kafka-pactflow-example-provider-js-kafka.json");

  const opts = {
    pactUrls: [pactUrl],
    messageProviders: {
      'a product event update': providerWithMetadata(() => createEvent(new Product("42", "food", "pizza"), "UPDATED"), {
        kafka_topic: 'products',
      }),
    },
  };

  const p = new MessageProviderPact(opts);

  describe("product api publishes an event", () => {
    it("can generate messages for specified consumers", () => {
      return p.verify();
    });
  });
});
