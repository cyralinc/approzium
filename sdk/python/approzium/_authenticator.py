# needed to be able to import protos code
import sys
from datetime import datetime, timedelta
from itertools import count
from pathlib import Path

import grpc

from . import _postgres, _mysql
from ._iam import (
    assume_role,
    get_local_arn,
    obtain_claimed_arn,
    obtain_credentials,
    obtain_signed_get_caller_identity,
)

sys.path.append(str(Path(__file__).parent / "_protos"))  # isort:skip
import authenticator_pb2  # noqa: E402 isort:skip
import authenticator_pb2_grpc  # noqa: E402 isort:skip


class AuthClient(object):
    """This class represents a connection to an Approzium authenticator
    service. Instances of this class can be used as arguments to database
    drivers connect method to use for authentication.

    :param server_address: address (host:port) at which an authenticator
        service is listening.

    :type server_address: str
    :param iam_role: if an IAM role Amazon resource number (ARN) is provided,
        it will be assumed and its identity will be used for authentication.
        Otherwise, the default ``boto3`` session will be used as the identity.
    :type iam_role: str, optional
    """

    def __init__(self, server_address, iam_role=None):
        self.server_address = server_address
        self.authenticated = False
        self._counter = count(1)
        self.n_conns = 0
        self.iam_role = iam_role

        # Parse the claimed ARN once because it'll never change.
        # Parse the signed_gci at startup, and then we'll update
        # it every 5 minutes because it expires every 15.
        if iam_role is None:
            claimed_arn = get_local_arn()
            signed_gci = obtain_signed_get_caller_identity(None)
        else:
            response = assume_role(iam_role)
            credentials = obtain_credentials(response)
            claimed_arn = obtain_claimed_arn(response)
            signed_gci = obtain_signed_get_caller_identity(credentials)

        self.claimed_arn = claimed_arn
        self.signed_gci = signed_gci
        self.signed_gci_last_updated = datetime.utcnow()

    @property
    def attribution_info(self):
        """Provides a dictionary containing information about the current state
        of the AuthClient. Useful for logging.

        :rtype: dict

        **Return Structure**:
            * *authenticator_address* (*str*): address of authenticator service used
            * *iam_role* (*str*): IAM Amazon resource number (ARN) used as identity
            * *authenticated* (*bool*): whether the AuthClient was verified by the
                                        authenticator service.
            * *num_connections* (*int*): number of connections made through this
                                         AuthClient
        """
        info = {}
        info["authenticator_address"] = self.server_address
        info["iam_role"] = self.iam_role
        info["authenticated"] = self.authenticated
        info["num_connections"] = self.n_conns
        return info

    def _execute_request(self, request, getmethodname):
        # The presigned GetCallerIdentity call expires every 15 minutes.
        self._update_gci_if_needed()

        channel = grpc.insecure_channel(self.server_address)
        stub = authenticator_pb2_grpc.AuthenticatorStub(channel)
        # add authentication info
        request.authtype = authenticator_pb2.AWS
        request.client_language = authenticator_pb2.PYTHON
        request.awsauth.CopyFrom(
            authenticator_pb2.AWSAuth(
                signed_get_caller_identity=self.signed_gci,
                claimed_iam_arn=self.claimed_arn,
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


    def _get_mysql_hash(self, dbhost, dbport, dbuser, auth_type, auth_info):
        if auth_type == _mysql.MYSQLNativePassword:
            salt = auth_info
            if len(salt) != 20:
                raise Exception("salt not right size")
            request = authenticator_pb2.MYSQLSHA1HashRequest(
                dbhost=dbhost, dbuser=dbuser, dbport=dbport, salt=salt,
            )
            response = self._execute_request(request, "GetMYSQLSHA1Hash")
            return response.hash
        else:
            raise Exception('Unexpected authentication method')

    # The presigned GetCallerIdentity string expires every 15 minuts, so refresh it
    # after 5 minutes just to be safe.
    def _update_gci_if_needed(self):
        if datetime.utcnow() - self.signed_gci_last_updated < timedelta(minutes=5):
            return
        if self.iam_role is None:
            self.signed_gci = obtain_signed_get_caller_identity(None)
            self.signed_gci_last_updated = datetime.utcnow()
        else:
            response = assume_role(self.iam_role)
            credentials = obtain_credentials(response)
            self.signed_gci = obtain_signed_get_caller_identity(credentials)
            self.signed_gci_last_updated = datetime.utcnow()
