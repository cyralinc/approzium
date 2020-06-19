all: server

server:
	go build .
deps:
	protoc -I protos authenticator.proto --go_out=plugins=grpc:protos
