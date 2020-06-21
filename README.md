# Tuto gRPC

Inspired by the course of Stephane Maarek.

- Unary
- Server Streaming
- Client Streaming
- Bi-Directional Streaming
- Blog example: CRUD + download image

## Documentation

https://grpc.io/

## Error Handling

- https://grpc.io/docs/guides/error/
- http://avi.im/grpc-errors/

## Deadlines

- https://grpc.io/blog/deadlines/

## Authentication

- https://grpc.io/docs/guides/auth/
- https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-auth-support.md

## Reflection

- https://github.com/grpc/grpc-go/tree/master/reflection
- https://github.com/ktr0731/evans
Command line:
go get github.com/ktr0731/evans
evans -p 50051 -r
show package
show service
show message
desc SumRequest
package calculator
service CalculatorService
call Sum

## mongoDB

- https://github.com/mongodb/homebrew-brew
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
mkdir database/data/db
mongod --dbpath database/data/db
- Interface for mongoDB - Robo 3T
https://robomongo.org/download
- driver
go get go.mongodb.org/mongo-driver/mongo

## Real world example

- https://github.com/googleapis/googleapis/blob/master/google/pubsub/v1/pubsub.proto
- https://github.com/googleapis/googleapis/blob/master/google/spanner/v1/spanner.proto 

## gRPC gateway

https://github.com/grpc-ecosystem/grpc-gateway

gRPC to JSON proxy generator following the gRPC HTTP spec.

## GoGo - Alternative of golang/protobuf

Third party with extra performance.

## Run blog

- local
export Environment=local

go run blog/blog_server/*.go 

go run blog/blog_client/client.go 

evans -p 50051 -r

- with TLS
export Environment=somethingElse

go run blog/blog_server/*.go 
evans -p 50051 -r --tls --host=localhost --cacert ssl/ca.crt  