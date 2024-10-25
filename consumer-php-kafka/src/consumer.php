<?php
/** @var \Enqueue\RdKafka\RdKafkaContext $context */

use DI\Container;
use Enqueue\RdKafka\RdKafkaConnectionFactory;
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Slim\Factory\AppFactory;

require __DIR__ . '/../vendor/autoload.php';


$container = new Container();
AppFactory::setContainer($container);
$app = AppFactory::create();

$container->set('repository', function () {
    return new ProductRepository();
});

$container->set('controller', function ($container) {
    return new ProductController($container->get('repository'));
});

$app->get('/products/{id}', 'controller:getById');
$app->get('/product/{id}', 'controller:getById');
$app->get('/products', 'controller:getAll');

class ProductController
{
    private $repository;

    public function __construct($repository)
    {
        $this->repository = $repository;
    }

    public function getAll($request, $response, $args)
    {
        $products = $this->repository->fetchAll();
        $response->getBody()->write(json_encode($products));
        return $response->withHeader('Content-Type', 'application/json');
    }

    public function getById($request, $response, $args)
    {
        $product = $this->repository->getById($args['id']);
        if ($product) {
            $response->getBody()->write(json_encode($product));
            return $response->withHeader('Content-Type', 'application/json');
        } else {
            $response->getBody()->write(json_encode(['message' => 'Product not found']));
            return $response->withStatus(404)->withHeader('Content-Type', 'application/json');
        }
    }
}

class Product
{
    public $id;
    public $type;
    public $name;
    public $version;

    public function __construct($id, $type, $name, $version)
    {
        if (!$id || !$type || !$name || !$version) {
            throw new InvalidArgumentException('All fields are mandatory');
        }

        $this->id = $id;
        $this->type = $type;
        $this->name = $name;
        $this->version = $version;
    }
}

class ProductEvent extends Product
{
    public $event;

    public function __construct($id, $type, $name, $version, $event)
    {
        parent::__construct($id, $type, $name, $version);

        if (!$event) {
            throw new InvalidArgumentException('Event field is mandatory');
        }

        $this->event = $event;
    }
}

class ProductRepository
{
    private $products;

    public function __construct()
    {
        $this->products = [
            "09" => new Product("09", "CREDIT_CARD", "Gem Visa", "v1"),
            "10" => new Product("10", "CREDIT_CARD", "28 Degrees", "v1"),
            "11" => new Product("11", "PERSONAL_LOAN", "MyFlexiPay", "v2"),
        ];
    }

    public function fetchAll()
    {
        return array_values($this->products);
    }

    public function getById($id)
    {
        return $this->products[$id] ?? null;
    }

    public function insert($product)
    {
        $this->products[$product->id] = $product;
    }

    public function delete($product)
    {
        unset($this->products[$product->id]);
    }
}

function productHandler($product)
{
    global $container;
    $repository = $container['repository'];
    print_r($product);
    print($product->id);
    print($product->type);
    print($product->name);
    print($product->version);
    if ($product->event == "CREATED" || $product->event == "UPDATED") {
        $repository->insert(new Product($product->id, $product->type, $product->name, $product->version));
    } elseif ($product->event == "DELETED") {
        $repository->delete(new Product($product->id, $product->type, $product->name, $product->version));
    } else {
        throw new Exception("Unable to process event");
    }
}

$connectionFactory = new RdKafkaConnectionFactory();

use Enqueue\RdKafka\Serializer;
use Enqueue\RdKafka\RdKafkaMessage;

class KafkaSerializer implements Serializer
{
    public function toString(RdKafkaMessage $message): string
    {
        $json = json_encode([
            'body' => $message->getBody(),
            'properties' => $message->getProperties(),
            'headers' => $message->getHeaders(),
        ]);

        if (JSON_ERROR_NONE !== json_last_error()) {
            throw new \InvalidArgumentException(sprintf(
                'The malformed json given. Error %s and message %s',
                json_last_error(),
                json_last_error_msg()
            ));
        }

        return $json;
    }

    public function toMessage(string $string): RdKafkaMessage
    {
        $data = json_decode($string, true);
        if (JSON_ERROR_NONE !== json_last_error()) {
            throw new \InvalidArgumentException(sprintf(
                'The malformed json given. Error %s and message %s',
                json_last_error(),
                json_last_error_msg()
            ));
        }
        return new RdKafkaMessage($string, [], []);
    }
}

$context = $connectionFactory->createContext();
$context->setSerializer(new KafkaSerializer());

$productsQueue = $context->createQueue('products');
$consumer = $context->createConsumer($productsQueue);

$message = $consumer->receive();
if ($message) {
    print($message->getBody());
    $productData = json_decode($message->getBody(), true);
    $productData = new ProductEvent(
        $productData['id'],
        $productData['type'],
        $productData['name'],
        $productData['version'],
        $productData['event']
    );
    try {
        productHandler($productData);
        $consumer->acknowledge($message);
    } catch (Exception $e) {
        error_log('unable to handle incoming message: ' . $e->getMessage());
    }
}

$app->run();
?>