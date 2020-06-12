import psycopg2
import select
import logging
import struct
from sys import getsizeof
import warnings
from ._psycopg2_ctypes import (
    libpq_PQstatus,
    libpq_PQsslInUse,
    libpq_PQgetssl,
    libpq_PQsetnonblocking,
    libssl_SSL_read,
    libssl_SSL_write,
    set_connection_sync,
    read_msg,
    write_msg,
    write_to_conn,
    set_debug
)
from .authenticator import get_hash
from .misc import read_int32_from_bytes
import approzium
import pyximport
pyximport.install(language_level=3)
from .pg_scram import SCRAMAuthentication


logger = logging.getLogger(__name__)

# Postgres protocol constants
# derived from PGsource/src/include/libpq/pgcomm.h
AUTH_REQ_MD5 = 5
AUTH_REQ_SASL = 10

pgconnect = psycopg2.connect



def read_auth(pgconn):
    # request many more bytes than necessary. if connection is at the
    # right stage, only the right number of bytes will be received
    msg_type, msg = read_msg(pgconn)
    auth_type = read_int32_from_bytes(msg, 0)
    if msg_type != b'R':
        raise Exception("Authentication message not received")
    if auth_type == AUTH_REQ_MD5:
        salt = msg[-4:]
        return auth_type, bytes(salt)
    elif auth_type == AUTH_REQ_SASL:
        if not msg[4:].startswith(b'SCRAM-SHA-256'):
            raise Exception("Server requested an unsupported SASL authentication method")
        auth = SCRAMAuthentication(b'SCRAM-SHA-256')
        dbuser = pgconn.get_dsn_parameters()["user"]
        client_first = auth.create_client_first_message(dbuser)
        select.select([], [pgconn.fileno()], [])
        write_msg(pgconn, b'p', client_first)
        select.select([pgconn.fileno()], [], [])
        resp_type, server_first = read_msg(pgconn)
        if resp_type != b'R':
            raise Exception('Error received unexpected response', server_first)
        # the part that is relevant is the part that starts with r=
        auth.parse_server_first_message(server_first[4:])
        return auth_type, auth
    else:
        raise Exception("Unidentified authentication method")


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


def send_hash(pgconn, auth_type, hash):
    if auth_type == AUTH_REQ_MD5:
        write_msg(pgconn, b'p', b'md5'+hash.encode('ascii')+b'\0')
    elif auth_type == AUTH_REQ_SASL:
        client_final, auth = hash
        write_msg(pgconn, b'p', client_final)
        select.select([pgconn.fileno()], [], [])
        resp_type, server_final = read_msg(pgconn)
        if resp_type != b'R':
            raise Exception('Error received unexpected response', server_first)
        if not auth.verify_server_final_message(server_final):
            raise Exception('Error bad server signature')


def construct_approzium_conn(base, is_sync):
    if not base:
        base = psycopg2.extensions.connection

    class ApproziumConn(base):
        CONNECTION_AWAITING_RESPONSE = 4

        def __init__(self, *args, **kwargs):
            # can safely do so because real async value was caught earlier in our connect method
            logger.debug("ApproziumConn __init__")
            kwargs.pop("async", None)
            kwargs.pop("async_", None)
            super().__init__(*args, **kwargs, async_=1)
            if self.dsn is None:
                # connection is uninitalized due to an error
                return
            if logger.getEffectiveLevel() <= logging.DEBUG:
                set_debug(self)
            self._salt = None
            self._auth_type = None
            self._hash_sent = False
            if is_sync:
                wait(self)
                set_connection_sync(self)
                self.autocommit = False

        def poll(self):
            status = libpq_PQstatus(self.pgconn_ptr)
            if status == self.CONNECTION_AWAITING_RESPONSE and not self._salt:
                logging.debug("reading salt")
                self._auth_type, self._salt = read_auth(self)
                return psycopg2.extensions.POLL_WRITE
            elif self._salt and not self._hash_sent:
                logging.debug("sending hash")
                dbhost = self.get_dsn_parameters()["host"]
                dbuser = self.get_dsn_parameters()["user"]
                hash = get_hash(
                    dbhost, dbuser, self._auth_type, self._salt, approzium.authenticator_addr
                )
                send_hash(self, self._auth_type, hash)
                self._hash_sent = True
                return psycopg2.extensions.POLL_WRITE
            else:
                logging.debug("normal poll")
                return super().poll()

    return ApproziumConn


def connect(dsn=None, connection_factory=None, cursor_factory=None, **kwargs):
    is_sync = True
    if kwargs.get("async", False):
        is_sync = False
    if kwargs.get("async_", False):
        is_sync = False
    # construct our approzium factory class on top of given connection factory class
    factory = construct_approzium_conn(connection_factory, is_sync)
    return pgconnect(dsn, factory, cursor_factory, **kwargs)
