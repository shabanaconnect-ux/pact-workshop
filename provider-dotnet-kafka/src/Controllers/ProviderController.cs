using Microsoft.AspNetCore.Mvc;

namespace Products.Controllers;

public class ProductsController : Controller
{
    private ProductRepository _repository;

    // This would usually be from a Repository/Data Store

    public ProductsController()
    {
        _repository = ProductRepository.GetInstance();
    }

    [HttpPost]
    [Route("/products")]
    public IActionResult Create([FromBody] Product product)
    {
        _repository.AddProduct(product);
        return new StatusCodeResult(201);
    }

    [HttpPut]
    [Route("/products/{id}")]
    public IActionResult Update(string id, [FromBody] Product product)
    {
        if (product != null)
        {
            _repository.UpdateProduct(product);
        return new OkResult();
        }
        return new NotFoundResult();
    }

    [HttpDelete]
    [Route("/products/{id}")]
    public IActionResult Delete(string id, [FromBody] Product product)
    {
        if (product != null)
        {
            _repository.DeleteProduct(product);
            return new OkResult();
        }
        return new NotFoundResult();
    }
}
