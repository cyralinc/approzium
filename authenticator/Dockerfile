FROM golang:1.13

# enable GOMODULES
ENV GO111MODULE on

# Install grpc
# Install protoc-gen-go protoc-gen-grpc-gateway and protoc-gen-swagger
RUN go get -u github.com/golang/protobuf/protoc-gen-go@v1.3.2

# Install protobuf compiler
RUN apt-get update && apt-get install -y unzip
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip
RUN unzip protoc-3.7.0-linux-x86_64.zip -d protoc3
RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/* /usr/local/include/

WORKDIR /usr/src/dbauth/authenticator