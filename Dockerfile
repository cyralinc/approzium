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
  libpq-dev
# Use Pyenv to install multiple Python versions for testing
RUN apt-get install -y libssl-dev zlib1g-dev libbz2-dev \
libreadline-dev libsqlite3-dev curl llvm libncurses5-dev libncursesw5-dev \
xz-utils tk-dev libffi-dev liblzma-dev python-openssl
RUN git clone git://github.com/yyuu/pyenv.git $HOME/.pyenv
ENV PYENV_ROOT $HOME/.pyenv
ENV PATH $PYENV_ROOT/shims:$PYENV_ROOT/bin:$PATH
RUN pyenv install 3.5.9
RUN pyenv install 3.6.10
RUN pyenv install 3.7.7
RUN pyenv install 3.8.3
RUN pyenv global 3.7.7
RUN pip3 install --upgrade pip
RUN pip3 install pytest
# Install protobuf compiler for Python
RUN pip3 install grpcio-tools
# Python package dependencies. Repeated here for faster cached build time
RUN pip3 install psycopg2 boto3 grpcio
WORKDIR /usr/src/approzium/sdk
COPY sdk/ .
# Install Python SDK in editable mode
RUN pip3 install -e python
# Build Authenticator Go Binary
WORKDIR /usr/src/approzium/authenticator
COPY authenticator/ .
RUN --mount=type=cache,target=$GOPATH/pkg/mod go build

FROM alpine:latest AS build
WORKDIR /app/
COPY --from=dev /usr/src/approzium/authenticator/authenticator .
ENTRYPOINT ["./authenticator"]
