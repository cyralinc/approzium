# needed to be able to import protos code
import sys
from itertools import count
from pathlib import Path

import grpc

from . import _postgres
from .iam import (
    assume_role,
    get_local_arn,
    obtain_claimed_arn,
    obtain_credentials,
    obtain_signed_get_caller_identity,
)

sys.path.append(str(Path(__file__).parent / "protos"))  # isort:skip
import authenticator_pb2  # noqa: E402 isort:skip
import authenticator_pb2_grpc  # noqa: E402 isort:skip


class AuthClient(object):
    def __init__(self, server_address, iam_role=None):
        self.server_address = server_address
        self.iam_role = iam_role
        self.authenticated = False
        self._counter = count(1)
        self.n_conns = 0

    @property
    def attribution_info(self):
        info = {}
        info["authenticator_address"] = self.server_address
        info["iam_role"] = self.iam_role
        info["authenticated"] = self.authenticated
        info["num_connections"] = self.n_conns
        return info


    def _execute_request(self, request, getmethodname):
        if self.iam_role is None:
            claimed_arn = get_local_arn()
            signed_gci = obtain_signed_get_caller_identity(None)
        else:
            response = assume_role(self.iam_role)
            credentials = obtain_credentials(response)
            claimed_arn = obtain_claimed_arn(response)
            signed_gci = obtain_signed_get_caller_identity(credentials)

        channel = grpc.insecure_channel(self.server_address)
        stub = authenticator_pb2_grpc.AuthenticatorStub(channel)
        # add authentication info
        request.authtype = authenticator_pb2.AWS
        request.client_language = authenticator_pb2.PYTHON
        request.awsauth.CopyFrom(
            authenticator_pb2.AWSAuth(
                signed_get_caller_identity=signed_gci,
                claimed_iam_arn=claimed_arn,
            )
        )
        response = getattr(stub, getmethodname)(request)
        # if no exception is raised, request was successful
        self.authenticated = True
        self.n_conns = next(self._counter)
        return response

    def _get_pg2_hash(self, dbhost, dbport, dbuser, auth_type, auth_info):
        if auth_type == _postgres.AUTH_REQ_MD5:
            salt = auth_info
            if len(salt) != 4:
                raise Exception("salt not right size")
            request = authenticator_pb2.PGMD5HashRequest(
                dbhost=dbhost, dbuser=dbuser, dbport=dbport, salt=salt,
            )
            response = self._execute_request(request, "GetPGMD5Hash")
            return response.hash
        elif auth_type == _postgres.AUTH_REQ_SASL:
            auth = auth_info
            auth._generate_auth_msg()
            request = authenticator_pb2.PGSHA256HashRequest(
                dbhost=dbhost,
                dbport=dbport,
                dbuser=dbuser,
                salt=auth.password_salt,
                iterations=auth.password_iterations,
                authentication_msg=auth.authorization_message,
            )
            response = self._execute_request(request, "GetPGSHA256Hash")
            client_final = auth.create_client_final_message(response.cproof)
            auth.server_signature = response.sproof
            return client_final, auth
