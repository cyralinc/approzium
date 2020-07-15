import logging

import asyncpg
from asyncpg.connect_utils import _ConnectionParameters
from asyncpg.protocol import Protocol

import approzium

from .._postgres import PGAuthClient, construct_msg, parse_msg

logger = logging.getLogger(__name__)


def appz_authenticate(self, data):
    # This method is used to override the method
    # `async.protocol.Protocol.data_received` and initiate and finish an
    # Approzium connection sequence. If the protocol instance is already
    # authenticated, this method reverts to the original `data_received` method

    context = self.password["approzium_context"]
    if context["authclient"] is None:
        # happens only once per connection
        context["authclient"] = PGAuthClient(
            None,  # read method is set on each run
            lambda msg: context["transport"].write(msg),
            context["authenticator"],
            context["host"],
            context["port"],
            self.user,
        )
    if not context["authclient"].done:
        # data could contain multiple message, in which case, we want AuthClient to
        # consume only the first message, while asyncpg consumes the rest
        msg_type, msg_content = parse_msg(data)
        first_msg_bytes = construct_msg(msg_type, msg_content)
        read_bytes = lambda: first_msg_bytes  # noqa: E731
        rest_bytes = data[len(first_msg_bytes) :]
        context["authclient"].read_bytes = read_bytes
        next(context["authclient"])
        if len(rest_bytes):
            Protocol.data_received(self, rest_bytes)
    else:
        Protocol.data_received(self, data)


def new_connection_made(self, transport):
    # `self` is a `asyncpg.protocol.Protocol` instance.

    # checks if this is an approzium connection
    if isinstance(self.password, dict) and "approzium_context" in self.password:
        # store transport because otherwise it is a private variable because of
        # Cython
        self.password["approzium_context"]["transport"] = transport

        # override the protocol's `data_received` method
        def new_data_received(data):
            appz_authenticate(self, data)

        self.data_received = new_data_received
    return original_connection_made(self, transport)


original_connection_made = Protocol.connection_made
Protocol.connection_made = new_connection_made


def new__connect_addr(*args, **kwargs):
    connection_class = kwargs["connection_class"]
    if connection_class.__name__ == "ApproziumConnection":
        host, port = kwargs["addr"]
        approzium_context = {
            "authenticator": connection_class.authenticator,
            "transport": None,  # this is determined later by `connection_made`
            "authclient": None,
            "host": host,
            "port": str(port),
        }
        # trick: put context in password field, discarding its passed value
        conn_params_dict = kwargs["params"]._asdict()
        conn_params_dict["password"] = {"approzium_context": approzium_context}
        new_conn_params = _ConnectionParameters(**conn_params_dict)
        kwargs["params"] = new_conn_params
    return original__connect_addr(*args, **kwargs)


original__connect_addr = asyncpg.connect_utils._connect_addr
asyncpg.connect_utils._connect_addr = new__connect_addr


def construct_approzium_connection(authenticator):
    # instantiate a fresh instance of the class on each connection
    # this instance is used to pass information that is needed during
    # connection process (e.g.: the authenticator)
    class ApproziumConnection(asyncpg.connection.Connection):
        pass

    ApproziumConnection.authenticator = authenticator
    return ApproziumConnection


async def connect(*args, authenticator=None, **kwargs):
    """Creates a Asyncpg connection through Approzium authentication. Takes
    the same arguments as ``asyncpg.connect``, in addition to the
    authenticator argument.

    :param authenticator: AuthClient instance to be used for authentication. If
        not provided, the default AuthClient, if set, is used.
    :type authenticator: approzium.AuthClient, optional
    :raises: TypeError, if no AuthClient is given and no default one is set.
    :rtype: ``asyncpg.Connection``

    Example:

    .. code-block:: python

        >>> import approzium
        >>> import asyncio
        >>> from approzium.asyncpg import connect
        >>> auth = approzium.AuthClient("myauthenticator.com:6001", disable_tls=True)
        >>> async def run():
        ...     con = await connect(user='postgres', authenticator=auth)
        ...     # use the connection just like any other Asyncpg connection
        ...     types = await con.fetch('SELECT * FROM pg_type')
        ...     print(types)
        >>> asyncio.get_event_loop().run_until_complete(run())
    """
    if authenticator is None:
        authenticator = approzium.default_auth_client
    conn = await asyncpg.connect(
        *args, **kwargs, connection_class=construct_approzium_connection(authenticator)
    )
    conn.authenticator = authenticator
    return conn
