# needed to be able to import protos code
import json
import sys
from datetime import datetime, timedelta
from itertools import count
from pathlib import Path

import grpc

from . import _iam, _mysql, _postgres
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

    :param disable_tls: defaults to False. When False, https is used and a
        client_cert and client_key proving the client's identity must be
        provided. When True, http is used and no other TLS options must be
        set.
    :type disable_tls: bool, optional

    :param tls_config: the TLS config to use for encrypted communication.
    :type tls_config: TLSConfig, optional

    :param iam_role: if an IAM role Amazon resource number (ARN) is provided,
        it will be assumed and its identity will be used for authentication.
        Otherwise, the default ``boto3`` session will be used as the identity.
    :type iam_role: str, optional
    """

    def __init__(
        self, server_address, disable_tls=False, tls_config=None, iam_role=None,
    ):
        self.server_address = server_address

        if not disable_tls:
            if tls_config is None:
                raise ValueError("if tls is not disabled, tls config must be provided")
            if tls_config.client_cert is None or tls_config.client_key is None:
                raise ValueError(
                    "if tls is not disabled, "
                    "client_cert and client_key must be provided"
                )
            self.tls_config = TLSConfig(
                trusted_certs=tls_config.trusted_certs,
                client_cert=tls_config.client_cert,
                client_key=tls_config.client_key,
            )

        self.disable_tls = disable_tls
        self.authenticated = False
        self._counter = count(1)
        self.n_conns = 0

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
            - *authenticator_address* (*str*): address of authenticator service used
            - *iam_arn* (*str*): IAM Amazon resource number (ARN) used as identity
            - *authenticated* (*bool*): whether the AuthClient was verified by the
              authenticator service.
            - *num_connections* (*int*): number of connections made through this
              AuthClient
        """
        info = {}
        info["authenticator_address"] = self.server_address
        info["iam_arn"] = self.claimed_arn
        info["authenticated"] = self.authenticated
        info["num_connections"] = self.n_conns
        info.update(_iam.attribution_info())
        return info

    @property
    def attribution_info_json(self):
        """Provides the same attribution info returned by
        :func:`~AuthClient.attribution_info` as a JSON format string

        :rtype: str
        """
        info = self.attribution_info
        return json.dumps(info)

    def _execute_request(self, request, getmethodname, dbhost, dbport, dbuser):
        # The presigned GetCallerIdentity call expires every 15 minutes.
        self._update_gci_if_needed()

        if self.disable_tls:
            channel = grpc.insecure_channel(self.server_address)
        else:
            credentials = grpc.ssl_channel_credentials(
                root_certificates=_read_file(self.tls_config.trusted_certs),
                certificate_chain=_read_file(self.tls_config.client_cert),
                private_key=_read_file(self.tls_config.client_key),
            )
            channel = grpc.secure_channel(self.server_address, credentials)

        stub = authenticator_pb2_grpc.AuthenticatorStub(channel)
        # authentication info
        aws_identity = authenticator_pb2.AWSIdentity(
            signed_get_caller_identity=self.signed_gci,
            claimed_iam_arn=self.claimed_arn,
        )
        password_request = authenticator_pb2.PasswordRequest(
            aws=aws_identity,
            client_language=authenticator_pb2.PYTHON,
            dbhost=dbhost,
            dbport=dbport,
            dbuser=dbuser,
        )
        request.pwd_request.CopyFrom(password_request)
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
            request = authenticator_pb2.PGMD5HashRequest(salt=salt,)
            response = self._execute_request(
                request, "GetPGMD5Hash", dbhost, dbport, dbuser
            )
            return response.hash
        elif auth_type == _postgres.AUTH_REQ_SASL:
            auth = auth_info
            auth._generate_auth_msg()
            request = authenticator_pb2.PGSHA256HashRequest(
                salt=auth.password_salt,
                iterations=auth.password_iterations,
                authentication_msg=auth.authorization_message,
            )
            response = self._execute_request(
                request, "GetPGSHA256Hash", dbhost, dbport, dbuser
            )
            client_final = auth.create_client_final_message(response.cproof)
            auth.server_signature = response.sproof
            return client_final, auth

    def _get_mysql_hash(self, dbhost, dbport, dbuser, auth_type, auth_info):
        if auth_type == _mysql.MYSQLNativePassword:
            salt = auth_info
            if len(salt) != 20:
                raise Exception("salt not right size")
            request = authenticator_pb2.MYSQLSHA1HashRequest(salt=salt,)
            response = self._execute_request(
                request, "GetMYSQLSHA1Hash", dbhost, dbport, dbuser
            )
            return response.hash
        else:
            raise Exception("Unexpected authentication method")

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


class TLSConfig(object):
    """This class represents the TLS config to be used while communicating
    with Approzium. Its fields are further described here:
    https://grpc.github.io/grpc/python/grpc.html#create-client-credentials

    :param trusted_certs: the path to the root certificate(s) that must have
        issued the identity certificate used by Approzium's authentication
        server.
    :type trusted_certs: str, optional

    :param client_cert: this client's certificate, used for proving its
        identity, and used by the caller to encrypt communication with
        its public key
    :type client_cert: str, optional

    :param client_key: this client's key, used for decrypting incoming
        communication that was encrypted by callers using the client_cert's
        public key
    :type client_key: str, optional
    """

    def __init__(self, trusted_certs=None, client_cert=None, client_key=None):
        self.trusted_certs = trusted_certs
        self.client_cert = client_cert
        self.client_key = client_key


def _read_file(path_to_file):
    with open(path_to_file, "rb") as f:
        return f.read()
