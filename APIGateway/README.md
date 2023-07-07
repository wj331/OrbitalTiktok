Tiktok orbital project: API gateway.

Key Features:

1. Request handling and routing with the Hertz library.
2. Microservices/ RPC communication via Apache Binary Thrift protocol using the Kitex library. Utilizes Generic Call feature of Kitex.
3. Service discovery using Nacos.
4. Rate limiting, based on users IP address.
5. Basic authentication.
6. Response caching with go-cache and hertz.
7. Request Pooling: The API Gateway has been optimized to use a pool of clients instead of creating a new client for every request.


Benchmarking results:

- 10000 requests, mean time per request for a Post request under /cache/post endpoint: 0.088ms. 


Routes are registered for each type of Generic Call. Below are the registered endpoints:

* `/jsonservice/:method`: Handles JSON Thrift generic call.
* `/bizservice/:method`: Handles HTTP generic call.
* `/post`: Handles Binary generic call.

* `/jsonprotoservice/:method`: Handles JSON Protobuf generic call.

Each endpoint utilizes a middleware for rate limiting (based on IP address). Caching can also be configured to improve the performance of frequent identical requests.

## Configuration
You can configure the rate limiting and caching parameters:

* `MaxQPS` is the maximum number of requests per second allowed from a single IP address.
* `BurstSize` is the maximum number of events that can occur at once.
* `cacheExpiryTime` is the duration for which data is stored in the cache.

### Development
* To add a new service, extend the pools map with the service name as key and the generic client pool as value.
* To register a new route, follow the format of the existing RegisterRoute functions.
