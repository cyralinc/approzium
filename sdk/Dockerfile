FROM ubuntu:18.04
ENV HOME /root
RUN apt-get update
RUN apt-get install -y build-essential
RUN apt-get install -y git
RUN apt-get install -y libpq-dev
# dependencies for compiling postgresql
RUN apt-get install -y libreadline-dev libssl-dev bison flex

# copy over dbauth libpq sdk because it is needed to compile modified postgreql
ADD internal /usr/src/dbauth/sdk/internal

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
ENV PATH="/usr/local/pgsql/bin:${PATH}"
RUN apt-get install -y python3-dev
RUN apt-get install -y python3-pip
RUN pip3 install psycopg2

WORKDIR /usr/src/dbauth
