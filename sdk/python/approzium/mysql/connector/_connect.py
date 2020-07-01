import approzium
import mysql.connector
from contextlib import contextmanager
from mysql.connector import MySQLConnection, CMySQLConnection

from ..._mysql import get_auth_resp




@contextmanager
def _patch__do_auth(sql_connection_class):
    original__do_auth = sql_connection_class._do_auth

    def _do_auth(self, *args, **kwargs):
        if self._password.__class__.__name__ != "AuthClient":
            return

        def _auth_response(
            client_flags, username, password, database, auth_plugin, auth_data, ssl_enabled
        ):
            authenticator = password
            is_secure_connection = (
                client_flags & mysql.connector.constants.ClientFlag.SECURE_CONNECTION
            )
            auth_response = get_auth_resp(
                authenticator,
                host,
                str(port),
                username,
                auth_plugin,
                auth_data,
                is_secure_connection,
            )
            return auth_response

        host = self.server_host
        port = self.server_port
        self._protocol._auth_response = _auth_response

        res = original__do_auth(self, *args, **kwargs)
        return res
    try:
        sql_connection_class._do_auth = _do_auth
        yield
    finally:
        sql_connection_class._do_auth = original__do_auth


def connect(*args, authenticator=None, **kwargs):
    if authenticator is None:
        authenticator = approzium.default_auth_client
    if authenticator is None:
        raise TypeError("Auth client not specified and not default auth client is set")
    kwargs["password"] = authenticator
    with _patch__do_auth(MySQLConnection):
        conn = mysql.connector.connect(*args, **kwargs)
    return conn
