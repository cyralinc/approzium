build-env:
	docker build -t dbauth-golang .

run-env:
	docker run -it -v "$(PWD)":/usr/src/dbauth/authenticator --network dbauthnet --rm --name dbauth-golang dbauth-golang bash

deps:
	protoc -I messages authenticator.proto --go_out=plugins=grpc:messages