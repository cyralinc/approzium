import mysql.connector
from mysql.connector.connection import MySQLConnection
from ..._mysql import get_auth_resp


def _do_auth(self, *args, **kwargs):
    def _auth_response(client_flags, username, password, database,
                       auth_plugin, auth_data, ssl_enabled):
        authenticator = password
        is_secure_connection = client_flags & \
            mysql.connector.constants.ClientFlag.SECURE_CONNECTION
        auth_response = get_auth_resp(authenticator, host, str(port), username,
                                      auth_plugin, auth_data, is_secure_connection)
        return auth_response

    host = self.server_host
    port = self.server_port
    self._protocol._auth_response = _auth_response

    res =  original__do_auth(self, *args, **kwargs)
    return res

original__do_auth = MySQLConnection._do_auth
MySQLConnection._do_auth = _do_auth


def connect(*args, authenticator=None, **kwargs):
    kwargs['password'] = authenticator
    conn = mysql.connector.connect(*args, **kwargs)
    return conn
