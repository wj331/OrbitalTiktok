## Overview
1. Service Discovery using Nacos. Go-routine updates a list of valid RPC instances providing respective services every minute.
2. Request handling and routing. RPC communication through Generic Call feature of Kitex. Newly implemented JSON Protobuf Codec for generic call is integrated into this API Gateway.
3. Requests target service is determined based on its HTTP structure / endpoint
4. Efficient Mapping of services to it's respective IDL files. Supports version control (different versions of IDL for same service), maintaining back-compatibility.
5. Request Pooling: This API Gateway utilizes a mapping of pools of generic clients, created as per the service_configs.json file. This mapping has its key based on: `servicename_version`
6. Supports dynamic handling of changes in IDL's. Changes are captured on the fly; the mapping of pool of generic clients & valid RPC instances for services is updated without having to re-deploy the API gateway
7. Rate limiting, based on clients IP address.
8. Response caching for identical requests.
9. Quickly & seamlessly handles large volume of requests. 0.032ms mean time per request across 200,000 concurrent requests during benchmarking.

## Code structure 
Implementation for API Gateway lies in `apiGateway.go`

Config folder contains:
* json files for configuration of API Gateway
* json files for configurations of Services

Pkg folder contains:
* Generic Client handlers
* Route handlers
* Service handlers

Utils folder contains:
* implementations for Service Discovery of RPC servers
* implementations for Rate limiting
* implementations for caching of requests

## Usage
To run the API Gateway:
`go run .`

## Routes
Routes are registered for each type of Generic Call. Below are the registered endpoints:
* `/jsonservice/:method`: Handles JSON Thrift generic call.
* `/bizservice/:method`: Handles HTTP generic call.
* `/post`: Handles Binary generic call.

* `/jsonprotoservice/:method`: Handles JSON Protobuf generic call.

Each endpoint utilizes a middleware for rate limiting (based on IP address). Caching can also be configured to improve the performance of frequent identical requests.

## Development
* To add a new service, extend the pools map with the service name as key and the generic client pool as value.
* To register a new route, follow the format of the existing RegisterRoute functions.
