#include <iostream>
#include "dbauth_client.h"
#include <cstddef>
#include <grpcpp/grpcpp.h>
#include "authenticator.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;

using namespace dbauth::authenticator::messages;
using namespace std;

int dbauth_get_hashed_password(char *user, char *md5Salt, int saltlen, char *result_buffer) {
    char *salt_string = (char *) malloc(saltlen+1);
    if (salt_string == NULL) {
    	return 0;
    }
    memcpy(salt_string, md5Salt, saltlen);
    salt_string[saltlen] = '\0';
    AuthenticateRequest request;
    request.set_salt(salt_string);
    request.set_identity(user);
    AuthenticateResponse reply;
    auto channel = grpc::CreateChannel("authenticator:1234", grpc::InsecureChannelCredentials());
    auto stub_ = Authenticator::NewStub(channel);

    ClientContext context;

    Status status = stub_->Authenticate(&context, request, &reply);
    free(salt_string);

    if (status.ok()) {
        if (reply.status() != AuthenticateResponse_Status_SUCCESS) {
            cout << reply.message() << endl;
            return 0;
        }
        auto hash = reply.credentials().hashedpassword();
        strcpy(result_buffer, hash.c_str());
        return 1;
    } else {
        cout << status.error_code() << ": " << status.error_message()
             << endl;
        return 0;
    }
}
