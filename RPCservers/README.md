## Key Features
1. Request handling.
2. Service Registry using Nacos. RPC servers automatically register themselves upon startup.
3. Sends back a simple response to the API Gateway, which then gets routed back to client.

## Usage & Configurations
To build & run the servers:
`go run main.go`

You can configure the number of servers built by changing the for loop count.
