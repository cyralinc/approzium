PGDIR = /root/postgresql
LIBPQDIR = $(PGDIR)/src/interfaces/libpq

all: build

build:
	gcc -I $(LIBPQDIR) -I $(PGDIR)/src/include -L $(LIBPQDIR) -L $(PGDIR)/src/port/ -o main main.c -lpgport -lpq
build-env:
	docker build -t dbauth-dev .
run-env:
	docker run -it -v "$(PWD)":/usr/src/dbauth --rm --name dbauth-dev dbauth-dev bash
