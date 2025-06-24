const {
  MatchersV3,
  PactV4,
} = require("@pact-foundation/pact");
const productEventHandler = require('./product.handler')
const { like, regex } = MatchersV3;

const path = require("path");

describe("Kafka handler", () => {
  const pact = new PactV4({
    consumer: "pactflow-example-consumer-js-kafka",
    dir: path.resolve(process.cwd(), "pacts"),
    provider: "pactflow-example-provider-js-kafka",
    logLevel: process.env.PACT_LOG_LEVEL ?? "info",
  });

  describe("receive a product update", () => {
    it("accepts a product event",  () => {

      pact
        .addAsynchronousInteraction("accepts a product event")
        .expectsToReceive("a product event update")
        .withContents({
          id: like("some-uuid-1234-5678"),
          type: like("Product Range"),
          name: like("Some Product"),
          version: like("v1"),
          event: regex("^(CREATED|UPDATED|DELETED)$","UPDATED"),
        })
        // .withMetadata({
        //   "contentType": "application/json",
        //   "kafka_topic": "products",
        // })
        .executeTest(async (message) => {
          console.log(message)
          const result = await productEventHandler(message)
          console.log(result)
            expect(result instanceof Map).toBe(true);
            expect(result.has("some-uuid-1234-5678")).toBe(true);
            expect(result.get("some-uuid-1234-5678")).toEqual({
            id: "some-uuid-1234-5678",
            type: "Product Range",
            name: "Some Product",
            version: "v1"
            });
        });
        });
    });
  });
