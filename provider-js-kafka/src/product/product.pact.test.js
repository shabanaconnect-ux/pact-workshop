const { MessageProviderPact, providerWithMetadata } = require("@pact-foundation/pact");
const { Product } = require("./product");
const { createEvent } = require("./product.event");

describe("Message provider tests", () => {

  const opts = {
    pactUrls: [process.env.PACT_URL],
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
