import logging
import hashlib
import grpc
from pathlib import Path

# needed to be able to import protos code
import sys

sys.path.append(str(Path(__file__).parent / "protos"))
import authenticator_pb2_grpc
import authenticator_pb2


def get_hash(dbuser, salt, authenticator):
    # XXX: for now
    if len(salt) != 4:
        raise Exception("salt not right size")
    password = "password"
    logging.debug(f"computing hash")
    first_hash = hashlib.md5((password + dbuser).encode("ascii")).hexdigest()
    second_hash = hashlib.md5(first_hash.encode("ascii") + salt).hexdigest()
    return second_hash
