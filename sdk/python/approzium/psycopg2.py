import psycopg2
import select
import socket
from socket import ntohl, htonl
import struct
import logging
import warnings
from sys import getsizeof
from ctypes import (
    cdll,
    create_string_buffer,
    string_at,
    memmove,
    c_void_p,
    c_int,
    c_char_p,
)
from ctypes.util import find_library
from .socketfromfd import fromfd
from .authenticator import get_hash
from .misc import read_int32_from_bytes
import approzium


pgconnect = psycopg2.connect


libpq = cdll.LoadLibrary("libpq.so.5")
libssl = cdll.LoadLibrary(find_library("ssl"))


# setup ctypes functions
# necessary to avoid segfaults when using multiple Python threads
libpq_PQstatus = libpq.PQstatus
libpq_PQstatus.argtypes = [c_void_p]
libpq_PQstatus.restype = c_int

libpq_PQsslInUse = libpq.PQsslInUse
libpq_PQsslInUse.argtypes = [c_void_p]
libpq_PQsslInUse.restype = c_int

libpq_PQgetssl = libpq.PQgetssl
libpq_PQgetssl.argtypes = [c_void_p]
libpq_PQgetssl.restype = c_void_p

libpq_PQsetnonblocking = libpq.PQsetnonblocking
libpq_PQsetnonblocking.argtypes = [c_void_p, c_int]
libpq_PQsetnonblocking.restype = c_int

libssl_SSL_read = libssl.SSL_read
libssl_SSL_read.argtypes = [c_void_p, c_char_p, c_int]
libssl_SSL_read.restype = c_int

libssl_SSL_write = libssl.SSL_write
libssl_SSL_write.argtypes = [c_void_p, c_char_p, c_int]
libssl_SSL_write.restype = c_int


def set_connection_sync(pgconn):
    mem = bytearray(string_at(id(pgconn), getsizeof(pgconn)))
    sizeofint = struct.calcsize("@i")
    sizeoflong = struct.calcsize("@l")

    def addressofint(number, mem=mem):
        int_bytes = struct.pack("@i", number)
        return mem.find(int_bytes)

    def intataddress(address):
        return struct.unpack("@i", mem[address : address + sizeofint])[0]

    # as a check, we check server and protocol version numbers, which succeed
    # the async value in the psycopg connection struct
    server_version_addr = addressofint(pgconn.server_version)
    # check that there is only one match for that value
    assert (
        addressofint(pgconn.server_version, mem[server_version_addr + sizeofint :])
        == -1
    )
    protocol_address = server_version_addr - sizeofint
    protocol_version = intataddress(protocol_address)
    assert protocol_version == pgconn.protocol_version
    async_address = protocol_address - sizeoflong
    async_value = struct.unpack("@l", mem[async_address:protocol_address])[0]
    assert async_value == pgconn.async
    new_async_value = struct.pack("@l", 0)
    memmove(id(pgconn) + async_address, new_async_value, sizeoflong)
    assert pgconn.async == 0
    error = libpq_PQsetnonblocking(pgconn.pgconn_ptr, 0)
    assert error == 0


def read_salt(pgconn):
    # request many more bytes than necessary. if connection is at the
    # right stage, only the right number of bytes will be received
    NBYTES = 8096
    if libpq_PQsslInUse(pgconn.pgconn_ptr):
        buffer = bytearray(NBYTES)
        c_buffer = create_string_buffer(bytes(buffer), NBYTES)
        ssl_obj = libpq_PQgetssl(pgconn.pgconn_ptr)
        nread = libssl_SSL_read(ssl_obj, c_buffer, NBYTES)
        challenge = bytearray(c_buffer.raw[:nread])
    else:
        fd = pgconn.fileno()
        with warnings.catch_warnings():
            warnings.simplefilter("ignore", ResourceWarning)
            sock = fromfd(fd)
            challenge = sock.recv(NBYTES)
    assert len(challenge) == 13, "Challenge length is not correct"
    for index, byte in enumerate(challenge):
        msg_size = read_int32_from_bytes(challenge, index + 1)
        auth_type = read_int32_from_bytes(challenge, index + 5)
        AUTH_MD5 = 5
        if byte == ord("R") and msg_size == 12 and auth_type == AUTH_MD5:
            salt = challenge[index + 9 : index + 9 + 4]
            return bytes(salt)
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
    if libpq_PQsslInUse(pgconn.pgconn_ptr):
        ssl_obj = libpq_PQgetssl(pgconn.pgconn_ptr)
        c_buffer = create_string_buffer(msg, len(msg))
        n = libssl_SSL_write(ssl_obj, c_buffer, len(msg))
        if n != len(msg):
            raise ValueError("could not send response")
    else:
        with warnings.catch_warnings():
            warnings.simplefilter("ignore", ResourceWarning)
            sock = fromfd(pgconn.fileno(), keep_fd=True)
            sock.sendall(msg)


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
                logging.debug("reading salt")
                self._salt = read_salt(self)
                return psycopg2.extensions.POLL_WRITE
            elif self._salt and not self._hash_sent:
                logging.debug("sending hash")
                dbhost = self.get_dsn_parameters()["host"]
                dbuser = self.get_dsn_parameters()["user"]
                hash = get_hash(
                    dbhost, dbuser, self._salt, approzium.authenticator_addr
                )
                send_hash(self, hash)
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

    # if pgconn is None or pgconn.pgconn_ptr is None:
    #    # if connection is uninitialized, something is wrong so return as is
    #    return pgconn
