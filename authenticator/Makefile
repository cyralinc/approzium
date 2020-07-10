all: server

dev:
	go build && go install

test:
	go test -v -race -p 1 ./...

server:
	go build .

deps:
	cd server && protoc -I protos authenticator.proto --go_out=plugins=grpc:protos && cd -
	cd server && protoc -I protos health.proto --go_out=plugins=grpc:protos && cd -
