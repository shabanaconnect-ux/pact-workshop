const { Kafka, Partitioners } = require("kafkajs");
const { productFromJson } = require("./product");
const { createEvent } = require("./product.event");

const KAFKA_BROKER = process.env.KAFKA_BROKER || "localhost:9092";
const TOPIC = process.env.TOPIC || "products";

class ProductEventService {
  constructor() {
    this.kafka = new Kafka({
      clientId: "product-event-service",
      brokers: [KAFKA_BROKER],
    });
    this.producer = this.kafka.producer({ createPartitioner: Partitioners.LegacyPartitioner });

    this.connect();

    process.on("exit", () => {
      this.disconnect();
    });
    process.on("SIGINT", () => {
      this.disconnect();
      process.exit();
    });
    process.on("SIGTERM", () => {
      this.disconnect();
      process.exit();
    });
  }

  async connect() {
    await this.producer.connect();
  }

  async disconnect() {
    await this.producer.disconnect();
  }

  async create(event) {
    const product = productFromJson(event);
    await this.publish(createEvent(product, "CREATED"));
  }

  async update(event) {
    const product = productFromJson(event);
    await this.publish(createEvent(product, "UPDATED"));
  }

  async delete(event) {
    const product = productFromJson(event);
    await this.publish(createEvent(product, "DELETED"));
  }

  async publish(message) {
    const payload = {
      topic: TOPIC,
      messages: [{ value: JSON.stringify(message) }],
    };

    console.log("ProductEventService - sending message:", message);

    return this.producer.send(payload);
  }
}

module.exports = {
  ProductEventService,
  createEvent,
};
