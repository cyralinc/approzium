all: server

server:
	go build .
deps:
	cd server && protoc -I protos authenticator.proto --go_out=plugins=grpc:protos && cd -
