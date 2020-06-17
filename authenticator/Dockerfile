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

# Install Vault CLI for ease of testing Vault out
RUN wget https://releases.hashicorp.com/vault/1.4.2/vault_1.4.2_linux_amd64.zip
RUN unzip vault_1.4.2_linux_amd64.zip
RUN mv vault /usr/local/bin/
WORKDIR /usr/src/approzium/authenticator
