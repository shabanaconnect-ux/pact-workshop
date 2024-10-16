using System.Threading.Tasks;
using PactNet;
using PactNet.Output.Xunit;
using Xunit;
using Xunit.Abstractions;
using Match = PactNet.Matchers.Match;
using Products;

namespace Consumer.Tests
{
    public class ProductEventProcessorTests
    {
        private readonly ProductEventProcessor consumer;

        private readonly IMessagePactBuilderV4 pact;

        public ProductEventProcessorTests(ITestOutputHelper output)
        {
            consumer = new ProductEventProcessor();

            var config = new PactConfig
            {
                PactDir = "../../../pacts/",
                Outputters = new[]
                {
                    new XunitOutput(output)
                }
            };

            pact = Pact.V4("Product API", "Event API", config).WithMessageInteractions();
        }

        [Fact]
        public async Task ProductEventHandler_ProductCreated_HandlesMessage()
        {
            await pact
                      .ExpectsToReceive("an event indicating that an order has been created")
                      .WithMetadata("kafka_topic","products")
                      .WithJsonContent(new
                      {
                          id = Match.Type("some-uuid-1234-5678"),
                          type = Match.Type("Product Range"),
                          name = Match.Type("Some Product"),
                          version = Match.Type("v1"),
                          @event = Match.Regex("UPDATED","^(CREATED|UPDATED|DELETED)$")
                      })
                      .VerifyAsync<ProductEventProcessor.ProductEvent>(consumer.ProductEventHandler);
        }
    }
}