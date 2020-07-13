import logging

from ._authenticator import AuthClient, TLSConfig

default_auth_client = None
"""
Set this variable to an instance of :class:`~approzium.AuthClient` to set it as
the default auth client to be used for connections.
"""

logger = logging.getLogger(__name__)
log_format = "[%(filename)s:%(lineno)s - %(funcName)10s() ] %(message)s"
formatter = logging.Formatter(log_format)
ch = logging.StreamHandler()
ch.setFormatter(formatter)
logger.addHandler(ch)


__all__ = ["default_auth_client", "AuthClient", "TLSConfig"]
