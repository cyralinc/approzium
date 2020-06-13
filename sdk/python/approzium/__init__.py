import logging
from .authenticator import Authenticator

authenticator_addr = "authenticator:1234"
iam_role = None

logger = logging.getLogger(__name__)
ch = logging.StreamHandler()
formatter = logging.Formatter("[%(filename)s:%(lineno)s - %(funcName)10s() ] %(message)s")
ch.setFormatter(formatter)
logger.addHandler(ch)


default_authenticator = None
def set_default_authenticator(new_authenticator):
    global default_authenticator
    default_authenticator = new_authenticator
