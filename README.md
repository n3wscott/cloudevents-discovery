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


---
Downstream demo:

We will build this graph up backwards and watch the flow of services.

```
Port:      :8282 <-- :8181 <-- :8080 
Adds:      a,b,c     c*,d       x,y,z

<-- denotes "pulls from"
```

Start `:8080`:
```
PORT=8080 DISCOVERY_DOWNSTREAM=http://localhost:8181 DISCOVERY_SERVICES_FILE=config/discovery/xyz.yaml go run ./cmd/server/
```

Start a watch on the lead
``
watch curl localhost:8080/services
``

Or fancy:

```
watch 'curl localhost:8080/services | jq ".[].description"'
```

Start `:8181`:

```
PORT=8181 DISCOVERY_DOWNSTREAM=http://localhost:8282 DISCOVERY_SERVICES_FILE=config/discovery/cd.yaml go run ./cmd/server/
```

Start `:8282`:

```
PORT=8282 DISCOVERY_SERVICES_FILE=config/discovery/abc.yaml go run ./cmd/server/
```

or to make a ring:

```
PORT=8282 DISCOVERY_DOWNSTREAM=http://localhost:8080 DISCOVERY_SERVICES_FILE=config/discovery/abc.yaml go run ./cmd/server/
```
