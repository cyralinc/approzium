MY_PATH = $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
DBAUTHDIR= $(realpath $(MK_PATH)../)

build-env:
	docker build -t dbauth-golang .

run-env:
	docker run -it -v "$(DBAUTHDIR)":/usr/src/dbauth/ --network dbauthnet --rm --name dbauth-golang dbauth-golang bash

deps:
	protoc -I messages authenticator.proto --go_out=plugins=grpc:messages
