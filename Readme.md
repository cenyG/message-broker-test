Simpler Message Broker
----------------------

###Generate proto files:
```bash
protoc --go_out=plugins=grpc:. proto/broker.proto
```

###Start evans for test:
```bash
evans ./proto/broker.proto -p 8080
```

###Start Docker with examples:
```bash
docker-compose up --build
```

###Integration tests:
```bash
go test broker/it
```