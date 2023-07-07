# Orbital API Gateway
This is an implementation of an API gateway that is a single point of entry into a system, built using cloudwego's Hertz and Kitex libraries. It routes requests to appropriate RPC servers simulating microservices, and employs caching and rate limiting strategies.

## Getting Started (API Gateway)
Follow the steps below to set up the API Gateway HTTP server.

### Prerequisites
* Golang version 1.14 or newer.
* A proper GOPATH environment.
* [Kitex installed](https://www.cloudwego.io/docs/kitex/getting-started/)

### Installation
Clone the repository to your local machine: <br>
`git clone https://github.com/simbayippy/OrbitalxTiktok.git`

Navigate into the cloned repository:
`cd OrbitalxTiktok/APIGateway`

Then, build and run the API Gateway:
`go run .`

## Usage
This API gateway serves as the single entry point into the system, dealing with HTTP requests from the clients, translating them into the appropriate RPC request formats, and forwarding these to the corresponding RPC servers. The response is then returned back to the client.

It supports the following Generic Calls of Kitex:
* JSON Mapping (thrift) Generic Call
* HTTP Mapping Generic Call
* Binary Generic Call

As well as an additional, newly implemented Generic Call:
* JSON (protobuf) Generic Call

## Getting Started (RPC servers)
Follow the steps below to set up the API Gateway HTTP server.

Navigate into the cloned repository:
`cd OrbitalxTiktok/RPCservers`

Currently, 4 types of RPC servers have been built to handle the different Generic Calls:
* JSON Thrift RPC
* HTTP RPC
* Binary RPC

and an additional:
* JSON proto RPC