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

## Deliverables accomplished
The main deliverable for this project was utilizing of the Thrift codec for Generic Call. Implementation for this lies under:
* `pkg/routes/routeJSONThrift.go` 
* `pkg/genericClients/genericJSONThriftClient.go`

Additional assignment was also accomplished. Implementation for the jsonproto_codec can be found [here](https://github.com/simbayippy/kitex)

Implementation for this new codec of in this API Gateway lies under:
* `pkg/routes/routeJSONProto.go` 
* `pkg/genericClients/genericJSONProtoClient.go`

## Usage
Set-up your configurations in `configs` folder

Run the API Gateway:
`go run .`

To test the API Gateway, send the curl command: <br>
`curl -X POST -H "Content-Type: application/json" -d @APIGateway/benchmark/postDataProto.json http://127.0.0.1:8080/JSONProtoService/v1.0.0/Test`

You should obtain the response: <br>
`"{\"response\":\"hello\"}"`

To run your own tests
1. add your thrift/protobuf file into the respective folders
2. add its service into the `configs/service_configs.json` file
3. include the generic client to be used and created through `GenericClientType`
4. create its own route handler. Follow the existing implementations
5. create your own RPC servers to handle this new service

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

Protobuf/Thrift folders
* stores the respective IDL files