# Pact Message workshop

## Introduction

This workshop is aimed at demonstrating core features and benefits of contract testing message based systems with Pact.


Modern distributed architectures are increasingly integrated in a decoupled, asynchronous fashion. Message queues such as ActiveMQ, RabbitMQ, SNS, SQS, Kafka and Kinesis are common, often integrated via small and frequent numbers of microservices (e.g. lambda). These sorts of interactions are referred to as "message pacts".

There are some minor differences between how Pact works in these cases when compared to the HTTP use case. Pact supports messages by abstracting away the protocol and specific queuing technology (such as Kafka) and focusses on the messages passing between them..

This workshop should take from 1 to 2 hours, depending on how deep you want to go into each topic.

**Workshop outline**:

- [step 1: **create consumer(subscriber)**](https://github.com/pact-foundation/pact-workshop-js/tree/step1#step-1---simple-consumer-calling-provider): Create our consumer before the Provider API even exists
- [step 2: **create consumer pact test**](https://github.com/pact-foundation/pact-workshop-js/tree/step3#step-3---pact-to-the-rescue): Write a Pact test for our consumer
- [step 3: **create provider(publisher)**](https://github.com/pact-foundation/pact-workshop-js/tree/step3#step-3---pact-to-the-rescue): Create our provider
- [step 4: **create provider pact test**](https://github.com/pact-foundation/pact-workshop-js/tree/step3#step-3---pact-to-the-rescue): Write a Pact test to verify our provider codebase
- [step 5: **extend consumer**](https://github.com/pact-foundation/pact-workshop-js/tree/step5#step-5---back-to-the-client-we-go): Extend our consumer

_NOTE: Each step is tied to, and must be run within, a git branch, allowing you to progress through each stage incrementally._

_EG: Move to step 2:_

_`git checkout step2`_

_`npm install`_

<hr/>

## Learning objectives

If running this as a team workshop format, you may want to take a look through the [learning objectives](./LEARNING.md).

## Requirements

[Docker](https://www.docker.com)

[Docker Compose](https://docs.docker.com/compose/install/)

[Node + NPM](https://nodejs.org/en/)

## Scenario

There are two components in scope for our workshop.

1. Product Catalog API (Consumer/Subscriber). a simple HTTP service that exposes product information as a REST API, which is fed events from an Event API on the `product` topic.
2. Product Event (Provider). A service that publishes product events to a Kafka stream on the `product` topic.
uct.

We will only be focusing on the messages sent via Kafka and not the HTTP service, you can follow the HTTP workshop, to see Pact in action against the HTTP endpoints.

## Step 1 - Create Consumer (subscriber)
## Step 2 - Create Consumer Pact test
## Step 3 - Create Provider (publisher)
## Step 4 - Create Provider Pact test
