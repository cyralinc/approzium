import logging
import select

import approzium
import psycopg2

from ..misc import read_int32_from_bytes
from ._psycopg2_ctypes import (
    ensure_compatible_ssl,
    libpq_PQstatus,
    read_msg,
    set_connection_sync,
    set_debug,
    write_msg,
)
from .._postgres import parse_auth_msg, AUTH_REQ_MD5, AUTH_REQ_SASL, PGAuthClient
from .._postgres.scram import SCRAMAuthentication

logger = logging.getLogger(__name__)


pgconnect = psycopg2.connect


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


def construct_approzium_conn(base, is_sync, authenticator):
    if not base:
        base = psycopg2.extensions.connection

    class ApproziumConn(base):
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
            dbhost = self.get_dsn_parameters()["host"]
            dbport = self.get_dsn_parameters()["port"]
            dbuser = self.get_dsn_parameters()["user"]
            self._pgauthclient = PGAuthClient(
                lambda: read_msg(self),
                lambda msg: write_msg(self, msg),
                authenticator,
                dbhost,
                dbport,
                dbuser
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
            if status == self.CONNECTION_AWAITING_RESPONSE and \
                not self._pgauthclient.done:
                next(self._pgauthclient)
                return psycopg2.extensions.POLL_WRITE
            else:
                logging.debug("normal poll")
                return super().poll()

    return ApproziumConn


def connect(
    dsn=None, connection_factory=None, cursor_factory=None, authenticator=None, **kwargs
):
    is_sync = True
    if kwargs.get("async", False):
        is_sync = False
    if kwargs.get("async_", False):
        is_sync = False
    if authenticator is None:
        authenticator = approzium.default_authenticator
    if authenticator is None:
        raise Exception("Authenticator not specified")
    # construct our approzium factory class on top of given connection factory class
    factory = construct_approzium_conn(connection_factory, is_sync, authenticator)
    return pgconnect(dsn, factory, cursor_factory, **kwargs)
