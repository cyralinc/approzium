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


class Authenticator(object):
    def __init__(self, address, iam_role=None):
        self.address = address
        self.iam_role = iam_role

def get_hash(dbhost, dbuser, auth_type, auth_info, authenticator):
    signed_gci = obtain_signed_get_caller_identity(authenticator.iam_role)
    channel = grpc.insecure_channel(authenticator.address)
    stub = authenticator_pb2_grpc.AuthenticatorStub(channel)

    if auth_type == approzium.psycopg2.AUTH_REQ_MD5:
        salt = auth_info
        if len(salt) != 4:
            raise Exception("salt not right size")
        request = authenticator_pb2.PGMD5HashRequest(
            signed_get_caller_identity=signed_gci,
            claimed_iam_arn=authenticator.iam_role,
            dbhost=dbhost,
            dbuser=dbuser,
            salt=salt,
        )
        response = stub.GetPGMD5Hash(request)
        return response.hash
    elif auth_type == approzium.psycopg2.AUTH_REQ_SASL:
        auth = auth_info
        request = authenticator_pb2.PGSHA256HashRequest(
            signed_get_caller_identity=signed_gci,
            claimed_iam_arn=authenticator.iam_role,
            dbhost=dbhost,
            dbuser=dbuser,
            salt=auth.password_salt,
            iterations=auth.password_iterations
        )
        response = stub.GetPGSHA256Hash(request)
        salted_password = response.spassword
        client_final = auth.create_client_final_message(salted_password)
        return client_final, auth
