using System;
using System.Text.Json;
using System.Threading.Tasks;
using Confluent.Kafka;

namespace Products;
public class KafkaProducer
{
    private readonly ProducerConfig _config;

    public KafkaProducer(string bootstrapServers = "localhost:9092")
    {
        _config = new ProducerConfig { BootstrapServers = bootstrapServers };
    }

    public async Task ProduceProductEventAsync(ProductEvent productEvent)
    {
        using var producer = new ProducerBuilder<Null, string>(_config).Build();
        try
        {
            var message = new Message<Null, string> { Value = JsonSerializer.Serialize(productEvent) };
            var deliveryResult = await producer.ProduceAsync("products", message);
            Console.WriteLine($"Delivered '{deliveryResult.Value}' to '{deliveryResult.TopicPartitionOffset}'");
        }
        catch (ProduceException<Null, string> e)
        {
            Console.WriteLine($"Delivery failed: {e.Error.Reason}");
        }
    }
}