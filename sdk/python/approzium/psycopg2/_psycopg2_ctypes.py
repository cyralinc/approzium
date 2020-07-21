import logging
import os
import select
import struct
import subprocess
from ctypes import (
    CDLL,
    c_char_p,
    c_int,
    c_void_p,
    cdll,
    create_string_buffer,
    memmove,
    string_at,
)
from ctypes.util import find_library
from os import path
from sys import getsizeof

logger = logging.getLogger(__name__)

libpq = cdll.LoadLibrary(find_library("pq"))


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


def stdout(command):
    return subprocess.run(command, capture_output=True).stdout.decode("utf-8")


def ssl_supported():
    out = stdout(["pg_config", "--configure"])
    return "--with-openssl" in out


def possible_library_files(name):
    return [
        "lib%s.dylib" % name,
        "%s.dylib" % name,
        "%s.framework/%s" % (name, name),
        "lib%s.so" % name,
    ]


def setup_ssl():
    sslpath = ""
    # try to find OpenSSL path in `pg_config`'s LDFLAGS
    out = stdout(["pg_config", "--ldflags"])
    for lib in out.split(" "):
        if "openssl" in lib:
            ssldir = lib.split("-L")[-1]
            # directory path is found, so search for actual file
            for filename in possible_library_files("ssl"):
                possible_sslpath = path.join(ssldir, filename)
                if path.exists(possible_sslpath):
                    sslpath = possible_sslpath
                    break
    # if none is found, use the SSL library that the system's dynamic linker finds
    if not sslpath:
        sslpath = find_library("ssl")
    global libssl
    global libssl_SSL_read
    global libssl_SSL_write
    libssl = cdll.LoadLibrary(sslpath)
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

    def ensure(check):
        if not check:
            raise Exception("Could not set connection to sync. Unidentified struct")

    # as a check, we check server and protocol version numbers, which succeed
    # the async value in the psycopg connection struct
    server_version_addr = addressofint(pgconn.server_version)
    # check that there is only one match for that value
    ensure(
        addressofint(pgconn.server_version, mem[server_version_addr + sizeofint :])
        == -1
    )
    protocol_address = server_version_addr - sizeofint
    protocol_version = intataddress(protocol_address)
    ensure(protocol_version == pgconn.protocol_version)
    async_address = protocol_address - sizeoflong
    async_value = struct.unpack("@l", mem[async_address:protocol_address])[0]
    ensure(async_value == pgconn.async_)
    new_async_value = struct.pack("@l", 0)
    memmove(id(pgconn) + async_address, new_async_value, sizeoflong)
    ensure(pgconn.async_ == 0)
    error = libpq_PQsetnonblocking(pgconn.pgconn_ptr, 0)
    ensure(error == 0)


def read_msg(pgconn):
    def read_bytes(n):
        if pgconn.info.ssl_in_use:
            buffer = bytearray(n)
            c_buffer = create_string_buffer(bytes(buffer), n)
            ssl_obj = libpq_PQgetssl(pgconn.pgconn_ptr)
            nread = -1
            while nread == -1:
                nread = libssl_SSL_read(ssl_obj, c_buffer, n)

            msg = bytes(c_buffer.raw[:nread])
            return msg
        else:
            fd = pgconn.fileno()
            return os.read(fd, n)

    select.select([pgconn.fileno()], [], [])
    msg_type = read_bytes(1)
    msg_size_b = read_bytes(4)
    msg_size = struct.unpack("!i", msg_size_b)[0]
    msg_content = read_bytes(msg_size - 4)
    logger.debug(f"read: {msg_type} {msg_content}")
    return msg_type + msg_size_b + msg_content


def write_msg(pgconn, msg):
    select.select([], [pgconn.fileno()], [])
    if libpq_PQsslInUse(pgconn.pgconn_ptr):
        ssl_obj = libpq_PQgetssl(pgconn.pgconn_ptr)
        c_buffer = create_string_buffer(msg, len(msg))
        n = libssl_SSL_write(ssl_obj, c_buffer, len(msg))
        if n != len(msg):
            raise ValueError("could not send response")
    else:
        os.write(pgconn.fileno(), msg)
    logger.debug(f"sent: {msg}")


def set_debug(conn):
    libc = CDLL(find_library("c"))
    stdout = c_void_p.in_dll(libc, "stdout")
    libpq.PQtrace(conn.pgconn_ptr, stdout)


def ensure_compatible_ssl(conn):
    if conn.info.ssl_attribute("library") != "OpenSSL":
        raise Exception("Unsupported SSL library")


if ssl_supported():
    setup_ssl()
