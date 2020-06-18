FROM golang:1.13

ENV HOME /root
ENV GOPATH /usr

# enable GOMODULES
ENV GO111MODULE on

# Install protoc-gen-go
RUN go get -u github.com/golang/protobuf/protoc-gen-go@v1.3.3

# Install protobuf compiler
RUN apt-get update && apt-get install -y unzip
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip
RUN unzip protoc-3.7.0-linux-x86_64.zip -d protoc3
RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/* /usr/local/include/

# gRPC CLI
WORKDIR $HOME/grpc
RUN git clone https://github.com/grpc/grpc.git .
RUN git submodule update --init
RUN apt-get install -y libgflags-dev cmake
RUN mkdir -p cmake/build && \
    cd cmake/build && \
    cmake -DgRPC_BUILD_TESTS=ON ../../ && \
    make grpc_cli && \
    cp grpc_cli /usr/bin/
WORKDIR /usr/src/approzium/authenticator
