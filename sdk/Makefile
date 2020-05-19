MY_PATH = $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
DBAUTHDIR= $(realpath $(MK_PATH)../)

AUTHMSGSDIR = authenticator/messages
PGDIR = ~/postgresql
LIBPQDIR = $(PGDIR)/src/interfaces/libpq

INCLUDES =
CC = gcc
IFLAGS =  -I $(LIBPQDIR) -I $(PGDIR)/src/include -I $(PGDIR)/src/include/libpq

CFLAGS = -g -Wall -Wextra -Werror -Wfatal-errors -Wno-unused-variable $(IFLAGS)
LDFLAGS = -L $(PGDIR)/src/port/ -L $(PGDIR)/src/interfaces/libpq/
LDLIBS = -lpgport -lpq

EXECUTABLES = main

all: $(EXECUTABLES) gen-proto

clean:
	rm -f $(EXECUTABLES) *.o $(AUTHMSGSDIR)/*.h $(AUTHMSGSDIR)/*.c

%.o:%.c $(INCLUDES)
	$(CC) $(CFLAGS) -c $<

main: main.o
	$(CC) -o main main.o $(LDFLAGS) $(LDLIBS)

gen-proto:
	cd .. && protoc --c_out=. $(AUTHMSGSDIR)/authenticator.proto

gen-libpqpatch:
	# following has to be done to ensure new files are included in patch
	(cd $(PGDIR) && git add -A && git diff --staged HEAD --no-prefix -u .) > libpqpatch.diff
	cd $(PGDIR) && git reset

build-env:
	docker build -t dbauth-sdk .

run-env:
	docker run -it -v "$(DBAUTHDIR)":/usr/src/dbauth --network dbauthnet --rm --name dbauth-sdk dbauth-sdk bash

docker-network:
	docker network create dbauthnet
