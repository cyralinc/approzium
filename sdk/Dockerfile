FROM ubuntu:18.04
ENV HOME /root
RUN apt-get update && apt-get install -y \
  build-essential \
  libpq-dev \
  iputils-ping
RUN apt-get install -y python3-dev python3-pip
RUN pip3 install --upgrade pip
# tools for generating Python protobuf files
RUN pip3 install grpcio-tools
# install sdk dependencies
COPY . /usr/src/approzium/sdk
WORKDIR /usr/src/approzium/sdk/python
RUN pip3 install -e .
