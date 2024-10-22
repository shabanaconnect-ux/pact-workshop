from flask import Flask, request, jsonify
from confluent_kafka import Producer
import json
import uuid

app = Flask(__name__)

# Kafka Producer Configuration
config = {
    'bootstrap.servers': 'localhost:9092',
    'acks': 'all'
}

producer = Producer(config)

class ProductRepository:
    @staticmethod
    def create_product(data):
        if "id" not in data:
            data["id"] = str(uuid.uuid4())
        ProductRepository._send_event('CREATED', data)
        return data

    @staticmethod
    def update_product(data):
        ProductRepository._send_event('UPDATED', data)
        return data

    @staticmethod
    def delete_product(data):
        ProductRepository._send_event('DELETED', data)
        return data

    @staticmethod
    def _send_event(event_type, data):
        product_event = {
            'event': event_type,
            'type': data["type"],
            'id': data["id"],
            'version': data["version"],
            'name': data["name"]
        }
        producer.produce('products', value=json.dumps(data).encode('utf-8'))
        producer.poll(10000)
        producer.flush()

class ProductController:
    @staticmethod
    def create():
        data = request.json
        product_id = ProductRepository.create_product(data)
        return jsonify({'product_id': product_id}), 201

    @staticmethod
    def update(product_id):
        data = request.json
        if ProductRepository.update_product(data):
            return jsonify({'message': 'Product updated'}), 200
        return jsonify({'message': 'Product not found'}), 404

    @staticmethod
    def delete(product_id):
        if ProductRepository.delete_product(data):
            return jsonify({'message': 'Product deleted'}), 200
        return jsonify({'message': 'Product not found'}), 404

# Routes
@app.route('/products', methods=['POST'])
def create_product():
    ProductController.create()
    return '', 200

@app.route('/products/<product_id>', methods=['PUT'])
def update_product(product_id):
    ProductController.update(product_id)
    return '', 201

@app.route('/products/<product_id>', methods=['DELETE'])
def delete_product(product_id):
    ProductController.delete(product_id)
    return '', 200

if __name__ == '__main__':
    app.run(port=8081, debug=True)