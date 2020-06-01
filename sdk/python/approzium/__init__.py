import grpc
from pathlib import Path
# needed to be able to import grpc code
import sys
sys.path.append(str(Path(__file__).parent / 'protos'))
import authenticator_pb2_grpc
import authenticator_pb2


def connect(dbname='db'):
    channel = grpc.insecure_channel('authenticator:1234')
    stub = authenticator_pb2_grpc.AuthenticatorStub(channel)
    request = authenticator_pb2.DBUserRequest(
            iam_role_arn='IAM_ARN',
            signed_get_caller_identity='SHEESH',
            dbname='DB1'
    )
    response = stub.GetDBUser(request)
    return 'sheesh'
