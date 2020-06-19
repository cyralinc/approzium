from os import environ
from approzium import Authenticator, set_default_authenticator
from approzium.psycopg2.pool import SimpleConnectionPool
import logging

logger = logging.getLogger("approzium")
# logger.setLevel(logging.DEBUG)


auth = Authenticator("authenticator:1234", iam_role=environ.get("TEST_IAM_ROLE"))
set_default_authenticator(auth)
conns = SimpleConnectionPool(50, 500, "")
conn = conns.getconn()
print("Connection Established")


def test():
    conn.cursor().execute("SELECT 1")


test()
