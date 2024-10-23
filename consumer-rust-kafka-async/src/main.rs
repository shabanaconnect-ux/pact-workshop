use actix_web::{web, App, HttpResponse, HttpServer, Responder};
use futures::StreamExt;
use rdkafka::config::ClientConfig;
use rdkafka::consumer::{Consumer, StreamConsumer};
use rdkafka::message::Message;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Mutex;

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Product {
    id: String,
    r#type: String,
    name: String,
    version: String,
}

#[derive(Serialize, Deserialize)]
pub struct ProductEvent {
    id: String,
    name: String,
    r#type: String,
    version: String,
    event: String,
}

pub struct AppState {
    products: Mutex<HashMap<String, Product>>,
}

async fn get_all(data: web::Data<AppState>) -> impl Responder {
    let products = data.products.lock().unwrap();
    let products: Vec<&Product> = products.values().collect();
    HttpResponse::Ok().json(products)
}

async fn get_by_id(data: web::Data<AppState>, product_id: web::Path<String>) -> impl Responder {
    let products = data.products.lock().unwrap();
    if let Some(product) = products.get(&product_id.into_inner()) {
        HttpResponse::Ok().json(product)
    } else {
        HttpResponse::NotFound().json("Product not found")
    }
}

pub fn product_event_processor(data: &web::Data<AppState>, payload: &[u8]) {
    let product_event: ProductEvent =
        serde_json::from_slice(payload).expect("Error deserializing product");
    let product = Product {
        id: product_event.id.clone(),
        r#type: product_event.r#type.clone(),
        name: product_event.name.clone(),
        version: product_event.version.clone(),
    };
    let mut products = data.products.lock().unwrap();
    match product_event.event.as_str() {
        "CREATED" | "UPDATED" => {
            products.insert(product_event.id.clone(), product);
        }
        "DELETED" => {
            products.remove(&product.id);
        }
        _ => {
            eprintln!("Unknown event type");
        }
    }
}

async fn kafka_consumer(data: web::Data<AppState>) {
    let consumer: StreamConsumer = ClientConfig::new()
        .set("group.id", "products-group")
        .set("bootstrap.servers", "localhost:9092")
        .create()
        .expect("Consumer creation failed");

    consumer
        .subscribe(&["products"])
        .expect("Can't subscribe to topic");

    let mut message_stream = consumer.stream();

    while let Some(message) = message_stream.next().await {
        match message {
            Ok(m) => {
                if let Some(payload) = m.payload() {
                    product_event_processor(&data, payload);
                }
            }
            Err(e) => eprintln!("Kafka error: {}", e),
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let products = Mutex::new(HashMap::new());
    let data = web::Data::new(AppState { products });

    // Start Kafka consumer
    let data_clone = data.clone();
    actix_rt::spawn(async move {
        kafka_consumer(data_clone).await;
    });

    HttpServer::new(move || {
        App::new()
            .app_data(data.clone())
            .route("/products", web::get().to(get_all))
            .route("/products/{id}", web::get().to(get_by_id))
            .route("/product/{id}", web::get().to(get_by_id))
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}

#[cfg(test)]
mod tests {

use expectest::{expect, prelude::be_some};
use pact_consumer::{matching_regex, prelude::*};
use serde_json::Value;
use crate::{product_event_processor, AppState};
use std::collections::HashMap;
use std::sync::Mutex;
use actix_web::web;
use expectest::matchers::be_equal_to;
#[test]
fn consumes_a_product_event_update_message() {
    // Define the Pact for the test (you can setup multiple interactions by chaining the given or message_interaction calls)
    // For messages we need to use the V4 Pact format.
    let mut pact_builder =
        // Define the message consumer and provider by name
        pact_consumer::builders::PactBuilder::new_v4("pactflow-example-consumer-rust-kafka", "pactflow-example-provider-rust-kafka");
    pact_builder

        // Adds an interaction given the message description and type.
        .message_interaction("a product event update", |mut i| {
            // Can set the test name (optional)
            i.test_name("consumes_a_product_event_update_message");
            // // defines a provider state. It is optional.
            // i.given("some state");
            // // defines a provider state with parameters. It is optional.
            // i.given_with_params("some state with params {param}",&json!({
            //     "param": "some param"
            //   }));
            // Set the contents of the message. Here we use a JSON pattern, so that matching rules are applied
            i.json_body(json_pattern!({
              "id": like!("some-uuid-1234-5678"),
              "type": like!("Product Range"),
              "name": like!("Some Product"),
              "version": like!("v1"),
              "event": matching_regex!("^(CREATED|UPDATED|DELETED)$","UPDATED")
            }));
            // Set any required metadata
            i.metadata("kafka_topic", "products");
            // Need to return the mutated interaction builder
            i
        });

    // Arrange. setup product database
    let products = Mutex::new(HashMap::new());
    let data = web::Data::new(AppState { products });
    
    // This will return each message configured with the Pact builder. We need to process them
    // with out message handler (it should be the one used to actually process your messages).
    for message in pact_builder.messages() {
        // Process the message here as it would if it came off the queue
        let message_bytes = message.contents.contents.value().unwrap();
        let kafka_topic = message.contents.metadata.get("kafka_topic");
        let _message: Value = serde_json::from_slice(&message_bytes).unwrap();

        // Send the message to our message processor
        product_event_processor(&data,&message_bytes);
        
        // assert of the state of our product database, after processing the message
        let products = data.products.lock().unwrap();
        let product = products.get("some-uuid-1234-5678").unwrap();
        println!("{:?}", product);
        expect!(product.id.as_str()).to(be_equal_to("some-uuid-1234-5678".to_string()));
        expect!(product.name.clone()).to(be_equal_to("Some Product".to_string()));
        expect!(product.r#type.clone()).to(be_equal_to("Product Range".to_string()));
        expect!(product.version.clone()).to(be_equal_to("v1".to_string()));

        // assert the correct topic is included in our message
        expect!(kafka_topic)
            .to(be_some().value("products"));
    }
}

}