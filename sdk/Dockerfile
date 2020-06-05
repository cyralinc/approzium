FROM ubuntu:18.04
ENV HOME /root
RUN apt-get update && apt-get install -y \
  build-essential \
  libpq-dev
RUN apt-get install -y python3-dev python3-pip
RUN pip3 install psycopg2 boto3
# protobuf tools for Python
RUN pip3 install grpcio grpcio-tools
WORKDIR /usr/src/dbauth/sdk
