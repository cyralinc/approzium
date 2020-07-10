# syntax=docker/dockerfile:1.0-experimental
FROM golang:1.13 AS dev
ENV HOME /root
ENV GOPATH /usr
# enable GOMODULES
ENV GO111MODULE on
ENV CGO_ENABLED 0
# Nice to haves for development
RUN apt-get update && apt-get install -y iputils-ping vim
# Install protoc-gen-go
RUN go get -u github.com/golang/protobuf/protoc-gen-go@v1.3.3
# Install protobuf compiler
RUN apt-get install -y unzip
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip
RUN unzip protoc-3.7.0-linux-x86_64.zip -d protoc3
RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/* /usr/local/include/
# Install Vault CLI for ease of testing Vault out
RUN wget https://releases.hashicorp.com/vault/1.4.2/vault_1.4.2_linux_amd64.zip
RUN unzip vault_1.4.2_linux_amd64.zip
RUN mv vault /usr/local/bin/
RUN apt-get install -y \
  build-essential \
  libpq-dev \
  python3.7 \
  python3.7-dev \
  python3-pip \
  python3-venv
RUN pip3 install poetry tox
WORKDIR /usr/src/approzium/sdk/python
COPY sdk/python .
RUN poetry run pip install -U pip setuptools
RUN poetry install --extras "sqllibs tracing"
# Build Authenticator Go Binary
WORKDIR /usr/src/approzium/authenticator
COPY authenticator/ .
RUN --mount=type=cache,target=$GOPATH/pkg/mod go build

FROM alpine:latest AS build
WORKDIR /app/
COPY --from=dev /usr/src/approzium/authenticator/authenticator .
ENTRYPOINT ["./authenticator"]
