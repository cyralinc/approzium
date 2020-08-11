import pymysql

import approzium

from .._misc import patch
from .._mysql import MYSQLNativePassword, get_auth_resp


def _scramble_native_password(context, salt):
    # use normal method if not Approzium connection
    if not isinstance(context, dict):
        return pymysql._auth.scramble_native_password(context, salt)
    return get_auth_resp(
        context["authenticator"],
        context["host"],
        str(context["port"]),
        context["user"],
        MYSQLNativePassword,
        salt,
    )


class ApproziumConnection(pymysql.connections.Connection):
    def __init__(self, *args, authenticator=None, **kwargs):
        """Creates a PyMySQL connection through Approzium authentication. Takes
        the same arguments as ``pymysql.connect``, in addition to the
        authenticator argument.

        :param authenticator: AuthClient instance to be used for authentication. If
            not provided, the default AuthClient, if set, is used.
        :type authenticator: approzium.AuthClient, optional
        :raises: TypeError, if no AuthClient is given and no default one is set.
        :rtype: ``pymysql.connections.Connection``

        Example:

        .. code-block:: python

            >>> import approzium
            >>> from approzium.pymysql import connect
            >>> auth = approzium.AuthClient("authenticatorhost:6001", disable_tls=True)
            >>> con = connect("host=DB.com dbname=mydb", authenticator=auth)
            >>> # use the connection just like any other PyMySQL connection
        """
        if authenticator is None:
            authenticator = approzium.default_auth_client
        if authenticator is None:
            raise TypeError(
                "Auth client not specified and no default auth client is set"
            )
        self.authenticator = authenticator
        return super(ApproziumConnection, self).__init__(*args, **kwargs)

    @patch(pymysql._auth, "scramble_native_password", _scramble_native_password)
    def _request_authentication(self):
        # store info needed for Approzium authentication in password
        self.password = {
            "authenticator": self.authenticator,
            "host": self.host,
            "port": self.port,
            "user": self.user,
        }
        return super(ApproziumConnection, self)._request_authentication()
