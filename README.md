# Orbital API Gateway
This is an implementation of an API gateway that is a single point of entry into a system, used for TikTok's architecture. It routes requests to appropriate microservices, employing caching and rate limiting strategies.

## Getting Started
Follow the steps below to set up the API Gateway HTTP server.

### Prerequisites
* Golang version 1.14 or newer.
* A proper GOPATH environment.

### Installation
Clone the repository to your local machine: <br>
`git clone https://github.com/simbayippy/OrbitalxTiktok/APIGateway.git`

Navigate into the cloned repository:
`cd APIGateway`

Then, build and run the application:
`go run .`

## Usage
This API gateway serves as the single entry point into the system, dealing with HTTP requests from the clients, translating them into the appropriate RPC requests, and forwarding these to the corresponding microservices. The response is then returned back to the client.

It supports the following Generic Calls of Kitex:

* JSON Mapping (thrift) Generic Call
* JSON Mapping (protobuf) Generic Call
* HTTP Mapping Generic Call
* Binary Generic Call


Routes are registered for each type of Generic Call. Below are the registered endpoints:

* `/jsonservice/:method`: Handles JSON Thrift generic call.
* `/jsonprotoservice/:method`: Handles JSON Protobuf generic call.
* `/bizservice/:method`: Handles HTTP generic call.
* `/post`: Handles Binary generic call.

Each endpoint utilizes a middleware for rate limiting (based on IP address). Caching can also be configured to improve the performance of frequent identical requests.

## Configuration
You can configure the rate limiting and caching parameters:

* `MaxQPS` is the maximum number of requests per second allowed from a single IP address.
* `BurstSize` is the maximum number of events that can occur at once.
* `cacheExpiryTime` is the duration for which data is stored in the cache.

### Development
* To add a new service, extend the pools map with the service name as key and the generic client pool as value.
* To register a new route, follow the format of the existing RegisterRoute functions.
