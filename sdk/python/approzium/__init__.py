import logging
from .authenticator import Authenticator


default_authenticator = None
logger = logging.getLogger(__name__)
log_format = "[%(filename)s:%(lineno)s - %(funcName)10s() ] %(message)s"
formatter = logging.Formatter(log_format)
ch = logging.StreamHandler()
ch.setFormatter(formatter)
logger.addHandler(ch)


def set_default_authenticator(new_authenticator):
    global default_authenticator
    default_authenticator = new_authenticator


__all__ = ["set_default_authenticator", "Authenticator"]
