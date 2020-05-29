FROM ubuntu:18.04
ENV HOME /root
RUN apt-get update
RUN apt-get install -y build-essential
RUN apt-get install -y git
# install postgresql from source
#     dependencies
RUN apt-get install -y libreadline-dev libssl-dev bison flex
WORKDIR $HOME/postgresql
RUN git clone git://git.postgresql.org/git/postgresql.git .
RUN ./configure --with-openssl --without-zlib --enable-thread-safety
RUN make -C src/bin install && \
make -C src/include install && \
make -C src/interfaces install
RUN cd src/interfaces/libpq && ln -s libpq.so dbauth_libpq.so
ENV PATH="/usr/local/pgsql/bin:${PATH}"
ENV LD_LIBRARY_PATH="/usr/local/pgsql/lib/:${LD_LIBRARY_PATH}"
RUN apt-get install -y python3-dev
RUN apt-get install -y python3-pip
RUN pip3 install psycopg2
WORKDIR /usr/src/dbauth/sdk
