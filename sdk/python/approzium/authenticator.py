# needed to be able to import protos code
import sys
from pathlib import Path

import approzium
import grpc

from .iam import (
    assume_role,
    obtain_claimed_arn,
    obtain_credentials,
    obtain_signed_get_caller_identity,
)

sys.path.append(str(Path(__file__).parent / "protos"))  # isort:skip
import authenticator_pb2  # noqa: E402 isort:skip
import authenticator_pb2_grpc  # noqa: E402 isort:skip


class Authenticator(object):
    def __init__(self, address, iam_role=None):
        self.address = address
        self.iam_role = iam_role


def get_hash(dbhost, dbport, dbuser, auth_type, auth_info, authenticator):
    response = assume_role(authenticator.iam_role)
    credentials = obtain_credentials(response)
    signed_gci = obtain_signed_get_caller_identity(credentials)
    channel = grpc.insecure_channel(authenticator.address)
    stub = authenticator_pb2_grpc.AuthenticatorStub(channel)

    if auth_type == approzium.psycopg2.AUTH_REQ_MD5:
        salt = auth_info
        if len(salt) != 4:
            raise Exception("salt not right size")
        request = authenticator_pb2.PGMD5HashRequest(
            authdata={
                "auth_type": "aws",
                "client_language": "python",
                "signed_get_caller_identity": signed_gci,
                "claimed_iam_arn": obtain_claimed_arn(response),
            },
            dbhost=dbhost,
            dbuser=dbuser,
            dbport=dbport,
            salt=salt,
        )
        response = stub.GetPGMD5Hash(request)
        return response.hash
    elif auth_type == approzium.psycopg2.AUTH_REQ_SASL:
        auth = auth_info
        auth._generate_auth_msg()
        request = authenticator_pb2.PGSHA256HashRequest(
            authdata={
                "auth_type": "aws",
                "client_language": "python",
                "signed_get_caller_identity": signed_gci,
                "claimed_iam_arn": obtain_claimed_arn(response),
            },
            dbhost=dbhost,
            dbport=dbport,
            dbuser=dbuser,
            salt=auth.password_salt,
            iterations=auth.password_iterations,
            authentication_msg=auth.authorization_message,
        )
        response = stub.GetPGSHA256Hash(request)
        client_final = auth.create_client_final_message(response.cproof)
        auth.server_signature = response.sproof
        return client_final, auth
