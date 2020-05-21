FROM ubuntu:18.04
ENV HOME /root
RUN apt-get update
RUN apt-get install -y build-essential
RUN apt-get install -y git
RUN apt-get install -y libpq-dev
# dependencies for compiling postgresql
RUN apt-get install -y libreadline-dev libssl-dev bison flex

# Install protobuf compiler
RUN apt-get install -y unzip wget
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip
RUN unzip protoc-3.7.0-linux-x86_64.zip -d protoc3
RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/* /usr/local/include/

# install protobuf gRPC
# copied from https://hub.docker.com/r/grpc/cxx/dockerfile 
RUN apt-get update && apt-get install -y \
    build-essential autoconf git pkg-config \
    automake libtool curl make g++ unzip \
    && apt-get clean

# install protobuf first, then grpc
ENV GRPC_RELEASE_TAG v1.10.x
RUN git clone -b ${GRPC_RELEASE_TAG} https://github.com/grpc/grpc /var/local/git/grpc && \
            cd /var/local/git/grpc && \
    git submodule update --init && \
    echo "--- installing protobuf ---" && \
    cd third_party/protobuf && \
    ./autogen.sh && ./configure --enable-shared && \
    make -j$(nproc) && make -j$(nproc) check && make install && make clean && ldconfig && \
    echo "--- installing grpc ---" && \
    cd /var/local/git/grpc && \
    make -j$(nproc) && make install && make clean && ldconfig

# copy over dbauth libpq sdk because it is needed to compile modified postgreql
ADD internal /usr/src/dbauth/sdk/internal
WORKDIR /usr/src/dbauth/sdk/internal
RUN make

# install postgresql from source
WORKDIR $HOME/postgresql
RUN git clone git://git.postgresql.org/git/postgresql.git .
# apply patch
COPY libpqpatch.diff .
RUN patch -p0 < libpqpatch.diff
RUN ./configure --with-openssl --without-zlib --enable-thread-safety
RUN make -C src/bin install && \
make -C src/include install && \
make -C src/interfaces install
RUN cd src/interfaces/libpq && ln -s libpq.so dbauth_libpq.so
ENV PATH="/usr/local/pgsql/bin:${PATH}"
RUN apt-get install -y python3-dev
RUN apt-get install -y python3-pip
RUN pip3 install psycopg2

WORKDIR /usr/src/dbauth/sdk
