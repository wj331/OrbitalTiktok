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
