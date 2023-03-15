protoc --proto_path=. --go_out=paths=source_relative:gen/go trip.proto
protoc --proto_path=. --go-grpc_out=paths=source_relative:gen/go trip.proto
protoc --proto_path=. --grpc-gateway_out=paths=source_relative,grpc_api_configuration=trip.yaml:gen/go trip.proto