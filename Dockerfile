FROM ubuntu:18.04
ENV HOME /root
RUN apt-get update
RUN apt-get install -y build-essential
RUN apt-get install -y git
RUN apt-get install -y libpq-dev
# dependencies for compiling postgresql
RUN apt-get install -y libreadline-dev libssl-dev bison flex
RUN apt-get install -y libprotobuf-c-dev libprotoc-dev protobuf-compiler

# install protobuf-c from source
WORKDIR $HOME/protobuf-c
Run apt-get install -y pkg-config dh-autoreconf
RUN git clone https://github.com/protobuf-c/protobuf-c.git .
RUN ./autogen.sh
RUN ./configure
RUN make && make install

# install protobuf-c-rpc from source
WORKDIR $HOME/protobuf-c-rpc
RUN git clone https://github.com/protobuf-c/protobuf-c-rpc.git .
RUN ./autogen.sh
RUN ./configure
RUN make && make install

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
# there might be a better way but for now this is needed in order to make libpq
# available in dynamic linking
# RUN ln -s /root/postgresql/src/interfaces/libpq/libpq.so.5 /usr/lib/libpq.so.5
ENV PATH="/usr/local/pgsql/bin:${PATH}"
RUN apt-get install -y python3-dev
RUN apt-get install -y python3-pip
RUN pip3 install psycopg2

WORKDIR /usr/src/dbauth
