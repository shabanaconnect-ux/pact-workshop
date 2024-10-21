
using System.Collections.Generic;
namespace Products;

public sealed class ProductRepository
{
  private static readonly ProductRepository _instance = new ProductRepository();
  private List<Product> _products = new List<Product>();

  private readonly KafkaProducer _kafkaProducer = new KafkaProducer();

  // Explicit static constructor to tell C# compiler
  // not to mark type as beforefieldinit
  static ProductRepository()
  {
  }

  private ProductRepository()
  {
    this._products.Add(new Product("1", "food", "pancake", "1.0.0"));
    this._products.Add(new Product("2", "food", "sanwhich", "1.0.0"));
  }

  public static ProductRepository GetInstance()
  {
    return _instance;
  }

  public void AddProduct(Product product) {
    _kafkaProducer.ProduceProductEventAsync(new ProductEventProducer().CreateEvent(product,"CREATED")).GetAwaiter().GetResult();
  }

  public void UpdateProduct(Product product) {
    _kafkaProducer.ProduceProductEventAsync(new ProductEventProducer().CreateEvent(product,"UPDATED")).GetAwaiter().GetResult();
  }

  public void DeleteProduct(Product product) {
    _kafkaProducer.ProduceProductEventAsync(new ProductEventProducer().CreateEvent(product,"DELETED")).GetAwaiter().GetResult();
  }
}
