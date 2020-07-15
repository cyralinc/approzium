from .._misc import patch
from .._mysql import get_auth_resp, MYSQLNativePassword
import pymysql


def _scramble_native_password(context, salt):
    return get_auth_resp(
        context['authenticator'],
        context['host'],
        str(context['port']),
        context['user'],
        MYSQLNativePassword,
        salt,
    )


class ApproziumConnection(pymysql.connections.Connection):
    def __init__(self, *args, authenticator=None, **kwargs):
        self.authenticator = authenticator
        return super(ApproziumConnection, self).__init__(*args, **kwargs)

    @patch(pymysql._auth, 'scramble_native_password', _scramble_native_password)
    def _request_authentication(self):
        # store info needed for Approzium authentication in password
        self.password = {'authenticator': self.authenticator,
                         'host': self.host,
                         'port': self.port,
                         'user': self.user,
                         }
        return super(ApproziumConnection, self)._request_authentication()
