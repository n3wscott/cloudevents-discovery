# cloudevents-discovery
PoC to implement the proposed CloudEvents Discovery API in golang.

Run with:

```shell
go run ./cmd/server
```

Poke with:

```shell
curl localhost:8080/services
curl localhost:8080/services/cbdd62e8-c095-11ea-b3de-0242ac130004
curl localhost:8080/services?name=widgets

curl localhost:8080/types
curl localhost:8080/types/com.example.widge.delete
curl localhost:8080/types?matching=create
```
