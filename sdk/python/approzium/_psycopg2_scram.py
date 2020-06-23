import base64
import re

# try to import the secrets library from Python 3.6+ for the
# cryptographic token generator for generating nonces as part of SCRAM
# Otherwise fall back on os.urandom
try:
    from secrets import token_bytes as generate_token_bytes
except ImportError:
    from os import urandom as generate_token_bytes


class SCRAMAuthentication:
    AUTHENTICATION_METHODS = [b"SCRAM-SHA-256"]
    DEFAULT_CLIENT_NONCE_BYTES = 16  # 24
    REQUIREMENTS_CLIENT_FINAL_MESSAGE = ["client_channel_binding", "server_nonce"]
    REQUIREMENTS_CLIENT_PROOF = [
        "password_iterations",
        "password_salt",
        "server_first_message",
        "server_nonce",
    ]

    def __init__(self, authentication_method):
        self.authentication_method = authentication_method
        self.authorization_message = None
        # channel binding is turned off for the time being
        self.client_channel_binding = b"n,,"
        self.client_first_message_bare = None
        self.client_nonce = None
        self.client_proof = None
        self.password_salt = None
        self.password_iterations = None
        self.server_first_message = None
        self.server_key = None
        self.server_nonce = None

    def create_client_first_message(self, username):
        """Create the initial client message for SCRAM authentication"""
        self.client_nonce = self._generate_client_nonce(self.DEFAULT_CLIENT_NONCE_BYTES)
        # set the client first message bare here, as it's used in a later step
        self.client_first_message_bare = (
            b"n=" + username.encode("utf-8") + b",r=" + self.client_nonce
        )
        # put together the full message here
        msg = bytes()
        msg += self.authentication_method + b"\0"
        client_first_message = (
            self.client_channel_binding + self.client_first_message_bare
        )
        msg += (len(client_first_message)).to_bytes(
            4, byteorder="big"
        ) + client_first_message
        return msg

    def create_client_final_message(self, client_proof):
        """Create the final client message as part of SCRAM authentication"""
        if any(
            [
                getattr(self, val) is None
                for val in self.REQUIREMENTS_CLIENT_FINAL_MESSAGE
            ]
        ):
            raise Exception("you need values from server to generate a client proof")

        # generate the client proof
        msg = bytes()
        msg += (
            b"c="
            + base64.b64encode(self.client_channel_binding)
            + b",r="
            + self.server_nonce
            + b",p="
            + client_proof.encode("ascii")
        )
        return msg

    def parse_server_first_message(self, server_response):
        """Parse the response from the first message from the server"""
        self.server_first_message = server_response
        try:
            self.server_nonce = re.search(
                b"r=([^,]+),", self.server_first_message
            ).group(1)
        except IndexError:
            raise Exception("could not get nonce")
        if not self.server_nonce.startswith(self.client_nonce):
            raise Exception("invalid nonce")
        try:
            self.password_salt = re.search(
                b"s=([^,]+),", self.server_first_message
            ).group(1)
        except IndexError:
            raise Exception("could not get salt")
        try:
            self.password_iterations = int(
                re.search(b"i=(\d+),?", self.server_first_message).group(1)  # noqa:W605
            )
        except (IndexError, TypeError, ValueError):
            raise Exception("could not get iterations")

    def verify_server_final_message(self, server_final_message):
        """Verify the final message from the server"""
        try:
            server_signature = re.search(b"v=([^,]+)", server_final_message).group(1)
        except IndexError:
            raise Exception("could not get server signature")

        return server_signature == self.server_signature.encode("ascii")

    def _generate_client_nonce(self, num_bytes):
        token = generate_token_bytes(num_bytes)

        return base64.b64encode(token)

    def _generate_auth_msg(self):
        self.authorization_message = (
            self.client_first_message_bare
            + b","
            + self.server_first_message
            + b",c="
            + base64.b64encode(self.client_channel_binding)
            + b",r="
            + self.server_nonce
        )
