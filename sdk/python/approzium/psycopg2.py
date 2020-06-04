import psycopg2
import select
import socket
from socket import ntohl, htonl
import struct
import logging
from sys import getsizeof
from ctypes import cdll, create_string_buffer, string_at, memmove
from ctypes.util import find_library
from .socketfromfd import fromfd
from .authenticator import get_hash
from .misc import read_int32_from_bytes


libpq = cdll.LoadLibrary("libpq.so.5")
libssl = cdll.LoadLibrary(find_library("ssl"))


def set_connection_sync(pgconn):
    mem = bytearray(string_at(id(pgconn), getsizeof(pgconn)))
    sizeofint = struct.calcsize("@i")
    sizeoflong = struct.calcsize("@l")

    def addressofint(number):
        int_bytes = struct.pack("@i", number)
        return mem.find(int_bytes)

    def intataddress(address):
        return struct.unpack("@i", mem[address : address + sizeofint])[0]

    # as a check, we check server and protocol version numbers, which succeed
    # the async value in the psycopg connection struct
    server_version_addr = addressofint(pgconn.server_version)
    protocol_address = server_version_addr - sizeofint
    protocol_version = intataddress(protocol_address)
    assert protocol_version == pgconn.protocol_version
    async_address = protocol_address - sizeoflong
    async_value = struct.unpack("@l", mem[async_address:protocol_address])[0]
    assert async_value == pgconn.async
    new_async_value = struct.pack("@l", 0)
    memmove(id(pgconn) + async_address, new_async_value, sizeoflong)
    assert pgconn.async == 0
    error = libpq.PQsetnonblocking(pgconn.pgconn_ptr, 0)
    assert error == 0


def advance_connection(pgconn):
    state = pgconn.poll()
    status = libpq.PQstatus(pgconn.pgconn_ptr)
    fd = pgconn.fileno()
    if state == psycopg2.extensions.POLL_OK:
        pass
    elif state == psycopg2.extensions.POLL_WRITE:
        select.select([], [fd], [])
    elif state == psycopg2.extensions.POLL_READ:
        select.select([fd], [], [])
    else:
        raise psycopg2.OperationalError("poll() returned %s" % state)
    return state, status


def advance_until_challenge(pgconn):
    CONNECTION_AWAITING_RESPONSE = 4
    while True:
        state, status = advance_connection(pgconn)
        if status == CONNECTION_AWAITING_RESPONSE:
            nbytes = 8096
            if libpq.PQsslInUse(pgconn.pgconn_ptr):
                buffer = bytearray(nbytes)
                c_buffer = create_string_buffer(bytes(buffer), nbytes)
                ssl_obj = libpq.PQgetssl(pgconn.pgconn_ptr)
                nread = libssl.SSL_read(ssl_obj, c_buffer, nbytes)
                challenge = bytearray(c_buffer.raw[:nread])
            else:
                challenge = sock.recv(nbytes)
            for index, byte in enumerate(challenge):
                msg_size = read_int32_from_bytes(challenge, index + 1)
                auth_type = read_int32_from_bytes(challenge, index + 5)
                AUTH_MD5 = 5
                if byte == ord("R") and msg_size == 12 and auth_type == AUTH_MD5:
                    salt = challenge[index + 9 : index + 9 + 4]
                    return salt
        elif state == psycopg2.extensions.POLL_OK:
            raise Exception("Connection already established")


def advance_until_end(pgconn):
    while True:
        state, status = advance_connection(pgconn)
        if state == psycopg2.extensions.POLL_OK:
            return


def send_hash(pgconn, hash):
    sock = fromfd(pgconn.fileno(), keep_fd=True)
    msg_length = 40  # fixed number for MD5 hash result
    # XXX: following code assumes protocol 3
    msg = b"p"  # message type
    msg += struct.pack("!i", msg_length)
    msg += b"md5"
    msg += hash.encode("ascii")
    msg += b"\0"
    if libpq.PQsslInUse(pgconn.pgconn_ptr):
        ssl_obj = libpq.PQgetssl(pgconn.pgconn_ptr)
        c_buffer = create_string_buffer(msg, len(msg))
        n = libssl.SSL_write(ssl_obj, c_buffer, len(msg))
        if n != len(msg):
            raise ValueError("could not send response")
    else:
        sock.sendall(msg)


def parse_dsn_args(dsn, appz_args):
    psycopg_dsn = ""
    for d in dsn.split():
        param, value = d.split("=")
        if param in appz_args:
            if appz_args[param] is None:
                appz_args[param] = value
        else:
            psycopg_dsn += " " + d

    if appz_args["authenticator"] is None:
        raise ValueError("Authenticator not specified")
    return psycopg_dsn


def connect(dsn="", authenticator=None, **psycopgkwargs):
    appz_args = {"authenticator": authenticator}
    psycopg_dsn = parse_dsn_args(dsn, appz_args)
    pgconn = psycopg2.connect(psycopg_dsn, **psycopgkwargs, async=1)
    dbuser = pgconn.get_dsn_parameters()["user"]
    salt = advance_until_challenge(pgconn)
    hash = get_hash(dbuser, salt, appz_args["authenticator"])
    logging.debug(f"salt: {salt}, hash: {hash}")
    send_hash(pgconn, hash)
    advance_until_end(pgconn)
    if async == 0:
        set_connection_sync(pgconn)
    return pgconn
