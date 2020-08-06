all: server

dev:
	go build && go install

test:
	go test -v -race ./...

server:
	go build .

deps:
	cd server && protoc -I protos authenticator.proto --go_out=plugins=grpc:protos && cd -

bin:
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS) ui' sh -c "'$(CURDIR)/scripts/build.sh'"
