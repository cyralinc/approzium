import psycopg2
import select
import socket
import struct
import logging
from .socketfromfd import fromfd
from .authenticator import get_hash
from .misc import read_int32_from_bytes


def advance(pgconn, until_salt=False):
    sock = fromfd(pgconn.fileno(), keep_fd=True)
    state = pgconn.poll()
    poll_read_count = 0
    while True:
        fd = pgconn.fileno()
        if state == psycopg2.extensions.POLL_OK:
            return
        elif state == psycopg2.extensions.POLL_WRITE:
            select.select([], [fd], [])
        elif state == psycopg2.extensions.POLL_READ:
            select.select([fd], [], [])
            if until_salt and poll_read_count == 0:
                bytes = sock.recv(8096)
                for index, byte in enumerate(bytes):
                    msg_size = read_int32_from_bytes(bytes, index + 1)
                    auth_type = read_int32_from_bytes(bytes, index + 5)
                    AUTH_MD5 = 5
                    if byte == ord("R") and msg_size == 12 and auth_type == AUTH_MD5:
                        salt = bytes[index + 9 : index + 9 + 4]
                        return salt
            poll_read_count += 1
        else:
            raise psycopg2.OperationalError("poll() returned %s" % state)
        state = pgconn.poll()


def send_hash(pgconn, hash):
    sock = fromfd(pgconn.fileno(), keep_fd=True)
    msg_length = 40  # fixed number for MD5 hash result
    # XXX: following code assumes protocol 3
    msg = b"p"  # message type
    msg += struct.pack("!i", msg_length)
    msg += b"md5"
    msg += hash.encode("ascii")
    msg += b"\0"
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
    psycopgkwargs["sslmode"] = "disable"  # XXX: for now
    appz_args = {"authenticator": authenticator}
    psycopg_dsn = parse_dsn_args(dsn, appz_args)
    pgconn = psycopg2.connect(psycopg_dsn, **psycopgkwargs, async=True)
    dbuser = pgconn.get_dsn_parameters()["user"]
    salt = advance(pgconn, until_salt=True)
    hash = get_hash(dbuser, salt, appz_args["authenticator"])
    logging.debug(f"salt: {salt}, hash: {hash}")
    send_hash(pgconn, hash)
    advance(pgconn)
    return pgconn
