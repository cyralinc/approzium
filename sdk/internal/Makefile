all: build

build:
	g++ -c -o dbauth_libpq.o dbauth_libpq.cc
	ar rvs libdbauthlibpq.a dbauth_libpq.o

clean:
	rm -f *.o
