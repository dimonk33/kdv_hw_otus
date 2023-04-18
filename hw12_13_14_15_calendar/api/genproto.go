package api

//go:generate protoc -I .  --go_out=. --go-grpc_out=. --grpc-gateway_out ./gen --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true EventService.proto
