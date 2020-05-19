#include <iostream>
#include "dbauth_libpq.h"
#include <cstddef>
#include <grpcpp/grpcpp.h>
#include "authenticator.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;

using namespace dbauth::authenticator::messages;
using namespace std;

char *dbauth_get_hashed_password(char *md5Salt) {
    AuthenticateRequest request;
    request.set_salt(md5Salt);
    request.set_identity("diotim");
    AuthenticateResponse reply;
    auto channel = grpc::CreateChannel("authenticator:1234", grpc::InsecureChannelCredentials());
    auto stub_ = Authenticator::NewStub(channel);

    ClientContext context;

    Status status = stub_->Authenticate(&context, request, &reply);

    if (status.ok()) {
        cout << reply.message() << endl;
        return NULL;
    } else {
        cout << status.error_code() << ": " << status.error_message()
             << endl;
        return NULL;
    }
}
