// #[cfg(test)]
// mod tests {

// use expectest::{expect, prelude::be_some};
// use pact_consumer::{matching_regex, prelude::*};
// use serde_json::Value;
// use crate::{product_event_processor, AppState};
// use std::collections::HashMap;
// use std::sync::Mutex;
// use actix_web::web;
// #[test]
// fn consumes_a_product_event_update_message() {
//     // Define the Pact for the test (you can setup multiple interactions by chaining the given or message_interaction calls)
//     // For messages we need to use the V4 Pact format.
//     let mut pact_builder =
//         // Define the message consumer and provider by name
//         pact_consumer::builders::PactBuilder::new_v4("pactflow-example-consumer-rust-kafka", "pactflow-example-provider-rust-kafka");
//     pact_builder

//         // Adds an interaction given the message description and type.
//         .message_interaction("a product event update", |mut i| {
//             // Can set the test name (optional)
//             i.test_name("consumes_a_product_event_update_message");
//             // // defines a provider state. It is optional.
//             // i.given("some state");
//             // // defines a provider state with parameters. It is optional.
//             // i.given_with_params("some state with params {param}",&json!({
//             //     "param": "some param"
//             //   }));
//             // Set the contents of the message. Here we use a JSON pattern, so that matching rules are applied
//             i.json_body(json_pattern!({
//               "id": like!("some-uuid-1234-5678"),
//               "type": like!("Product Range"),
//               "name": like!("Some Product"),
//               "version": like!("v1"),
//               "event": matching_regex!("^(CREATED|UPDATED|DELETED)$","UPDATED")
//             }));
//             // Set any required metadata
//             i.metadata("kafka_topic", "products");
//             // Need to return the mutated interaction builder
//             i
//         });

//     // This will return each message configured with the Pact builder. We need to process them
//     // with out message handler (it should be the one used to actually process your messages).
//     for message in pact_builder.messages() {
//         let message_bytes = message.contents.contents.value().unwrap();
//         let kafka_topic = message.contents.metadata.get("kafka_topic");

//         // Process the message here as it would if it came off the queue
//         let message: Value = serde_json::from_slice(&message_bytes).unwrap();

//         let products = Mutex::new(HashMap::new());
//         let data = web::Data::new(AppState { products });
//         product_event_processor(data,message);
//         // Make some assertions on the processed value
//         expect!(message.as_object().unwrap().get("id"))
//             .to(be_some().value("some-uuid-1234-5678"));
//         expect!(kafka_topic)
//             .to(be_some().value("products"));
//     }
// }

// }