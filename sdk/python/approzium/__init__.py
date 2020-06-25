import logging

from .authenticator import AuthClient

default_auth_client = None
logger = logging.getLogger(__name__)
log_format = "[%(filename)s:%(lineno)s - %(funcName)10s() ] %(message)s"
formatter = logging.Formatter(log_format)
ch = logging.StreamHandler()
ch.setFormatter(formatter)
logger.addHandler(ch)


__all__ = ["default_auth_client", "AuthClient"]
