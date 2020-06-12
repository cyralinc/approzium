from ctypes import (
    cdll,
    CDLL,
    create_string_buffer,
    string_at,
    memmove,
    c_void_p,
    c_int,
    c_char_p,
)
from ctypes.util import find_library
import socket
import struct
import logging
from sys import getsizeof
import warnings
from .socketfromfd import fromfd


logger = logging.getLogger(__name__)

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


def read_from_conn(pgconn, nbytes, peek=True, justbytes=False):
    ### XXX: replace with conn.info.ssl_in_use
    if libpq_PQsslInUse(pgconn.pgconn_ptr) and not justbytes:
        buffer = bytearray(nbytes)
        c_buffer = create_string_buffer(bytes(buffer), nbytes)
        ssl_obj = libpq_PQgetssl(pgconn.pgconn_ptr)
        nread = libssl_SSL_read(ssl_obj, c_buffer, nbytes)
        msg = bytearray(c_buffer.raw[:nread])
    else:
        fd = pgconn.fileno()
        with warnings.catch_warnings():
            warnings.simplefilter("ignore", ResourceWarning)
            flags = [socket.MSG_PEEK] if peek else []
            sock = fromfd(fd)
            msg = sock.recv(nbytes, *flags)
    logger.debug(f'got: {msg}')
    return msg

def write_to_conn(pgconn, msg):
    if libpq_PQsslInUse(pgconn.pgconn_ptr):
        ssl_obj = libpq_PQgetssl(pgconn.pgconn_ptr)
        c_buffer = create_string_buffer(msg, len(msg))
        n = libssl_SSL_write(ssl_obj, c_buffer, len(msg))
        if n != len(msg):
            raise ValueError("could not send message")
    else:
        with warnings.catch_warnings():
            warnings.simplefilter("ignore", ResourceWarning)
            sock = fromfd(pgconn.fileno(), keep_fd=True)
            sock.sendall(msg)
    logger.debug(f'sent: {msg}')

def set_debug(conn):
    libc = CDLL(find_library('c'))
    stdout = c_void_p.in_dll(libc, 'stdout')
    libpq.PQtrace(conn.pgconn_ptr, stdout)
