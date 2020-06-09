import approzium
import logging
import hashlib
import grpc
from pathlib import Path

# needed to be able to import protos code
import sys

sys.path.append(str(Path(__file__).parent / "protos"))
import authenticator_pb2_grpc
import authenticator_pb2
from .iam import obtain_signed_get_caller_identity


def get_hash(dbhost, dbuser, salt, authenticator):
    if len(salt) != 4:
        raise Exception("salt not right size")
    iam_arn, signed_gci = obtain_signed_get_caller_identity()
    request = authenticator_pb2.PGMD5HashRequest(
        signed_get_caller_identity=signed_gci,
        claimed_iam_arn=iam_arn,
        dbhost=dbhost,
        dbuser=dbuser,
        salt=salt,
    )
    channel = grpc.insecure_channel(authenticator)
    stub = authenticator_pb2_grpc.AuthenticatorStub(channel)
    response = stub.GetPGMD5Hash(request)
    return response.hash
