FROM ubuntu:18.04
ENV HOME /root
RUN apt-get update
RUN apt-get install -y build-essential
RUN apt-get install -y libpq-dev
RUN apt-get install -y git
# dependencies for compiling postgresql
RUN apt-get install -y libreadline-dev libssl-dev bison flex

# install postgresql from source
WORKDIR $HOME/postgresql
RUN git clone git://git.postgresql.org/git/postgresql.git .
RUN ./configure --with-openssl --without-zlib
RUN make

RUN apt-get install -y libprotobuf-c-dev libprotoc-dev protobuf-compiler

# install protobuf-c from source
WORKDIR $HOME/protobuf-c
Run apt-get install -y pkg-config dh-autoreconf
RUN git clone https://github.com/protobuf-c/protobuf-c.git .
RUN ./autogen.sh
RUN ./configure
RUN make && make install

WORKDIR /usr/src/dbauth
