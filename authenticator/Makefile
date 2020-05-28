all: server

server:
	# hack in order to build server with flat file system
	-mv client.go client
	-go build .
	-mv client client.go
deps:
	protoc -I protos authenticator.proto --go_out=plugins=grpc:protos
