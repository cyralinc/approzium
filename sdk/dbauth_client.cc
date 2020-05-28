#include <iostream>
#include "dbauth_client.h"
#include <cstddef>
#include <grpcpp/grpcpp.h>
#include "authenticator.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;

using namespace dbauth::authenticator::protos;
using namespace std;

int dbauth_get_user(char *authenticatoraddress, char *identity, char **dbuser) {
    DBUserRequest request;
    request.set_identity(identity);
    DBUserResponse response;
    auto channel = grpc::CreateChannel(authenticatoraddress, grpc::InsecureChannelCredentials());
    auto stub_ = Authenticator::NewStub(channel);
    ClientContext context;
    Status status = stub_->GetDBUser(&context, request, &response);
    if (status.ok()) {
        *dbuser = (char *) malloc(response.dbuser().length()+1);
        strcpy(*dbuser, response.dbuser().c_str());
        return 1;
    } else {
        cout << status.error_code() << ": " << status.error_message()
             << endl;
        return 0;
    }
}
int dbauth_get_hash(char *authenticatoraddress, char *identity, char *md5salt, int saltlen, char *hashbuffer) {
    char *salt_string = (char *) malloc(saltlen+1);
    if (salt_string == NULL) {
    	return 0;
    }
    memcpy(salt_string, md5salt, saltlen);
    salt_string[saltlen] = '\0';
    DBHashRequest request;
    request.set_salt(salt_string);
    request.set_identity(identity);
    DBHashResponse reply;
    auto channel = grpc::CreateChannel(authenticatoraddress, grpc::InsecureChannelCredentials());
    auto stub_ = Authenticator::NewStub(channel);

    ClientContext context;

    Status status = stub_->GetDBHash(&context, request, &reply);
    free(salt_string);

    if (status.ok()) {
        auto hash = reply.hash();
        strcpy(hashbuffer, hash.c_str());
        return 1;
    } else {
        cout << status.error_code() << ": " << status.error_message()
             << endl;
        return 0;
    }
}
