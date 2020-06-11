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
    read_from_conn,
    write_to_conn,
)
from .authenticator import get_hash
from .misc import read_int32_from_bytes, redirect_socket_nowhere
import approzium


# Postgres protocol constants
# derived from PGsource/src/include/libpq/pgcomm.h
AUTH_REQ_MD5 = 5
AUTH_REQ_SASL = 10

pgconnect = psycopg2.connect


def _normal_poll(pgconn):
    '''Poll an Approzium type connection using just the Psycopg2 poll and not
    our custom poll function.
    '''
    super(type(pgconn), pgconn).poll()

def read_salt(pgconn):
    # request many more bytes than necessary. if connection is at the
    # right stage, only the right number of bytes will be received
    NBYTES = 8096
    # peek now and only remove data from socket once we are sure that this is
    # a supported authentication message.
    challenge = read_from_conn(pgconn, NBYTES, peek=True)
    msg_size = read_int32_from_bytes(challenge, 1)
    auth_type = read_int32_from_bytes(challenge, 5)
    if challenge[0] != ord("R"):
        raise Exception("Authentication message not received")
    if auth_type == AUTH_REQ_MD5 and msg_size == 12:
        salt = challenge[9 : 9 + 4]
        # remove bytes from socket and feed them to libpq through void socket
        read_from_conn(pgconn, NBYTES, peek=False)
        with redirect_socket_nowhere(pgconn.fileno(), feedit=challenge):
            try:
                _normal_poll(pgconn)
            except:
                pass
        return bytes(salt)
    elif auth_type == AUTH_REQ_SASL and msg_size == 23:
        if challenge[9:22] != b'SCRAM-SHA-256':
            raise Exception("Server requested an unsupported SASL authentication method")
        import pdb; pdb.set_trace()
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


def send_hash(pgconn, hash):
    msg_length = 40  # fixed number for MD5 hash result
    # XXX: following code assumes protocol 3
    msg = b"p"  # message type
    msg += struct.pack("!i", msg_length)
    msg += b"md5"
    msg += hash.encode("ascii")
    msg += b"\0"
    write_to_conn(pgconn, msg)


def construct_approzium_conn(base, is_sync):
    if not base:
        base = psycopg2.extensions.connection

    class ApproziumConn(base):
        CONNECTION_AWAITING_RESPONSE = 4

        def __init__(self, *args, **kwargs):
            # can safely do so because real async value was caught earlier in our connect method
            logging.debug("ApproziumConn __init__")
            kwargs.pop("async", None)
            kwargs.pop("async_", None)
            conn = super().__init__(*args, **kwargs, async_=1)
            if self.dsn is None:
                # connection is uninitalized due to an error
                return
            self._salt = None
            self._hash_sent = False
            if is_sync:
                wait(self)
                set_connection_sync(self)
                self.autocommit = False

        def poll(self):
            status = libpq_PQstatus(self.pgconn_ptr)
            if status == self.CONNECTION_AWAITING_RESPONSE and not self._salt:
                logger.debug("reading salt")
                self._salt = read_salt(self)
                return psycopg2.extensions.POLL_WRITE
            elif self._salt and not self._hash_sent:
                logger.debug("sending hash")
                dbhost = self.get_dsn_parameters()["host"]
                dbuser = self.get_dsn_parameters()["user"]
                hash = get_hash(
                    dbhost, dbuser, self._salt, approzium.authenticator_addr
                )
                send_hash(self, hash)
                self._hash_sent = True
                return psycopg2.extensions.POLL_WRITE
            else:
                logger.debug("normal poll")
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
