PGDIR = /root/postgresql
LIBPQDIR = $(PGDIR)/src/interfaces/libpq

EXECUTABLES = main

INCLUDES = dbauth.h
CC = gcc
IFLAGS =  -I $(LIBPQDIR) -I $(PGDIR)/src/include 

CFLAGS = -g -Wall -Wextra -Werror -Wfatal-errors -Wno-unused-variable $(IFLAGS)
LDFLAGS = -L $(PGDIR)/src/port/
LDLIBS = -lpgport -lpq


all: $(EXECUTABLES)

clean:
	rm -f $(EXECUTABLES) *.o

%.o:%.c $(INCLUDES)
	$(CC) $(CFLAGS) -c $<

main: main.o dbauth.o
	$(CC) $(LDFLAGS) -o main main.o dbauth.o $(LDLIBS)

build-env:
	docker build -t dbauth-dev .

run-env:
	docker run -it -v "$(PWD)":/usr/src/dbauth --rm --name dbauth-dev dbauth-dev bash

gen-proto:
	protoc --c_out:. authenticator/messages/authenticator.proto