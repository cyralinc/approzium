"""
Internal module that implements Approzium routines that are common to Postgres
libraries and the Postgres protocol.
"""
import logging
import struct

from .._misc import read_int32_from_bytes
from .scram import SCRAMAuthentication

logger = logging.getLogger(__name__)

# Postgres protocol constants
# derived from PGsource/src/include/libpq/pgcomm.h
AUTH_REQ_MD5 = 5
AUTH_REQ_SASL = 10


def parse_msg(msg):
    msg_type = msg[0]
    msg_size = struct.unpack("!i", msg[1 : 1 + 4])[0]
    msg_content = msg[5 : 5 + msg_size - 4]
    return msg_type, msg_content


def construct_msg(header, msg):
    if isinstance(header, int):
        header = header.to_bytes(1, "big")
    return header + struct.pack("!i", len(msg) + 4) + msg


def parse_auth_msg(msg):
    auth_type = read_int32_from_bytes(msg, 0)
    auth_info = {}
    if auth_type == AUTH_REQ_MD5:
        auth_info["salt"] = bytes(msg[-4:])
    elif auth_type == AUTH_REQ_SASL:
        if not msg[4:].startswith(b"SCRAM-SHA-256"):
            raise Exception("Server requested an unsupported SASL auth method")
        auth_info["scram"] = SCRAMAuthentication(b"SCRAM-SHA-256")
    return auth_type, auth_info


class PGAuthClient(object):
    """Class used to implement behaviour and state for an Approzium Postgres
    connection"""

    def __init__(self, read_bytes, write_bytes, authenticator, dbhost, dbport, dbuser):
        """When read_func is called with no arguments, a server message is
        received. When write_func is called with some bytes, they are written
        to the server.
        An instance should be instantiated when the next read_bytes call will
        return the first authentication request sent by the server."""
        self.read_bytes = read_bytes
        self.write_bytes = write_bytes
        self.authenticator = authenticator
        self.dbhost = dbhost
        self.dbport = dbport
        self.dbuser = dbuser
        self.done = False
        self.next_call = self.read_first_auth_msg

    def __next__(self):
        self.next_call()

    def read_first_auth_msg(self):
        msg_type, msg = parse_msg(self.read_bytes())
        if msg_type != ord("R"):
            raise Exception("Authentication message not received")
        auth_type, self.auth_info = parse_auth_msg(msg)
        if auth_type == AUTH_REQ_MD5:
            hash = self.authenticator._get_pg2_hash(
                self.dbhost,
                self.dbport,
                self.dbuser,
                AUTH_REQ_MD5,
                self.auth_info["salt"],
            )
            msg = construct_msg(b"p", b"md5" + hash.encode("ascii") + b"\0")
            self.write_bytes(msg)
            self.done = True
        elif auth_type == AUTH_REQ_SASL:
            scram_state = self.auth_info["scram"]
            client_first = scram_state.create_client_first_message(self.dbuser)
            msg = construct_msg(b"p", client_first)
            self.write_bytes(msg)
            self.next_call = self.scram_stage_2
        else:
            raise Exception("Unidentified authentication method")

    def scram_stage_2(self):
        msg_type, server_first = parse_msg(self.read_bytes())
        if msg_type != ord("R"):
            raise Exception("Error received unexpected response", server_first)
        # the part that is relevant is the part that starts with r=
        scram_state = self.auth_info["scram"]
        scram_state.parse_server_first_message(server_first[4:])
        client_final = self.authenticator._get_pg2_hash(
            self.dbhost,
            self.dbport,
            self.dbuser,
            AUTH_REQ_SASL,
            self.auth_info["scram"],
        )[0]
        msg = construct_msg(b"p", client_final)
        self.write_bytes(msg)
        self.next_call = self.scram_stage_3

    def scram_stage_3(self):
        msg_type, server_final = parse_msg(self.read_bytes())
        if msg_type != ord("R"):
            raise Exception("Error received unexpected response", server_final)
        scram_state = self.auth_info["scram"]
        if not scram_state.verify_server_final_message(server_final):
            raise Exception("Error bad server signature")
        self.done = True
