using System;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace Products
{
    public class ProductEventProcessor
    {
        private readonly ProductRepository _repository;

        public ProductEventProcessor()
        {
            _repository = ProductRepository.GetInstance();
        }

        public class ProductEvent(string id, string type, string name, string version, string @event) : Product(id, type, name, version)
        {
            [JsonPropertyName("event")]
            public string Event { get; set; } = @event;
        }

        public Task ProductEventHandler(ProductEvent productEvent)
        {
            Console.WriteLine($"Received product: {productEvent}");
            Console.WriteLine($"Received product event: {productEvent.Event}");
            Console.WriteLine($"Received product id: {productEvent.id}");

            if (productEvent.Event == "CREATED" || productEvent.Event == "UPDATED")
            {
                _repository.AddProduct(new Product(productEvent.id,productEvent.type,productEvent.name,productEvent.version));
            }
            else if (productEvent.Event == "DELETED")
            {
                _repository.RemoveProduct(productEvent.id);
            }
            else
            {
                throw new InvalidOperationException("Unable to process event");
            }

            return Task.CompletedTask;
        }
    }
}