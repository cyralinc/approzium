all: libdbauthlibpq.so main

CXX = g++
CPPFLAGS += `pkg-config --cflags protobuf grpc++`
CXXFLAGS += -std=c++11 -fPIC
SHEESH = -L/usr/local/lib `pkg-config --libs protobuf grpc++`\
		   -lgrpc++_reflection -ldl

libdbauthlibpq.so: authenticator.pb.o authenticator.grpc.pb.o dbauth_libpq.o
	$(CXX) -shared -o $@ $^ $(SHEESH)

main: main.o libdbauthlibpq.so
	$(CXX) main.o -L. -Wl,-rpath=. -ldbauthlibpq -o $@

deps:
	protoc --grpc_out=. --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` authenticator.proto --proto_path=../../authenticator/messages/
	protoc --cpp_out=. authenticator.proto --proto_path=../../authenticator/messages/

clean:
	rm -f *.o *.a
