import logging
import select

import psycopg2

import approzium

from .._postgres import PGAuthClient
from ._psycopg2_ctypes import (
    ensure_compatible_ssl,
    libpq_PQstatus,
    read_msg,
    set_connection_sync,
    set_debug,
    write_msg,
)

logger = logging.getLogger(__name__)


def wait(pgconn):
    while True:
        state = pgconn.poll()
        if state == psycopg2.extensions.POLL_OK:
            break
        elif state == psycopg2.extensions.POLL_WRITE:
            select.select([], [pgconn.fileno()], [])
        elif state == psycopg2.extensions.POLL_READ:
            select.select([pgconn.fileno()], [], [])
        else:
            raise psycopg2.OperationalError("poll() returned %s" % state)


def construct_approzium_conn(is_sync, authenticator):
    class ApproziumConn(psycopg2.extensions.connection):
        CONNECTION_AWAITING_RESPONSE = 4

        def __init__(self, *args, **kwargs):
            logger.debug("ApproziumConn __init__")
            kwargs.pop("async", None)
            kwargs.pop("async_", None)
            super().__init__(*args, **kwargs, async_=1)
            if self.dsn is None:
                # connection is uninitalized due to an error
                return
            if logger.getEffectiveLevel() <= logging.DEBUG:
                set_debug(self)
            self.authenticator = authenticator
            dbhost = self.get_dsn_parameters()["host"]
            dbport = self.get_dsn_parameters()["port"]
            dbuser = self.get_dsn_parameters()["user"]
            self._pgauthclient = PGAuthClient(
                lambda: read_msg(self),
                lambda msg: write_msg(self, msg),
                authenticator,
                dbhost,
                dbport,
                dbuser,
            )
            self._checked_ssl = False
            if is_sync:
                wait(self)
                set_connection_sync(self)
                self.autocommit = False

        def poll(self):
            status = libpq_PQstatus(self.pgconn_ptr)
            if self.info.ssl_in_use and not self._checked_ssl:
                ensure_compatible_ssl(self)
                logging.debug("checked ssl")
                self._checked_ssl = True
            if (
                status == self.CONNECTION_AWAITING_RESPONSE
                and not self._pgauthclient.done
            ):
                next(self._pgauthclient)
                return psycopg2.extensions.POLL_WRITE
            else:
                logging.debug("normal poll")
                return super().poll()

    return ApproziumConn


def connect(dsn=None, cursor_factory=None, authenticator=None, **kwargs):
    """Creates a Psycopg2 connection through Approzium authentication. Takes
    the same arguments as ``psycopg2.connect``, in addition to the
    authenticator argument.

    :param authenticator: AuthClient instance to be used for authentication. If
        not provided, the default AuthClient, if set, is used.
    :type authenticator: approzium.AuthClient, optional
    :raises: TypeError, if no AuthClient is given and no default one is set.
    :rtype: ``psycopg2.Connection``

    Example:

    .. code-block:: python

        >>> import approzium
        >>> from approzium.psycopg2 import connect
        >>> auth = approzium.AuthClient("myauthenticator.com:6001", disable_tls=True)
        >>> con = connect("host=DB.com dbname=mydb", authenticator=auth)  # no password!
        >>> # use the connection just like any other Psycopg2 connection

    .. warning::
        Currently, only `psycopg2` with dynamically linked `libpq` is
        supported. Thus, `psycopg2-binary` is not supported.
    """
    is_sync = True
    if kwargs.get("async", False):
        is_sync = False
    if kwargs.get("async_", False):
        is_sync = False
    if authenticator is None:
        authenticator = approzium.default_auth_client
    if authenticator is None:
        raise TypeError("Auth client not specified and not default auth client is set")
    factory = construct_approzium_conn(is_sync, authenticator)
    return psycopg2.connect(dsn, factory, cursor_factory, **kwargs)
