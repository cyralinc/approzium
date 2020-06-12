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
from . import psycopg2


def get_hash(dbhost, dbuser, auth_type, salt, authenticator):
    iam_arn, signed_gci = obtain_signed_get_caller_identity()
    channel = grpc.insecure_channel(authenticator)
    stub = authenticator_pb2_grpc.AuthenticatorStub(channel)

    if auth_type == psycopg2.AUTH_REQ_MD5:
        if len(salt) != 4:
            raise Exception("salt not right size")
        request = authenticator_pb2.PGMD5HashRequest(
            signed_get_caller_identity=signed_gci,
            claimed_iam_arn=iam_arn,
            dbhost=dbhost,
            dbuser=dbuser,
            salt=salt,
        )
        response = stub.GetPGMD5Hash(request)
        return response.hash
    elif auth_type == psycopg2.AUTH_REQ_SASL:
        auth = salt
        client_final = auth.create_client_final_message('password')
        return client_final, auth
