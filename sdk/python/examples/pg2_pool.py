from os import environ

from approzium import Authenticator
from approzium.psycopg2.pool import ThreadedConnectionPool

auth = Authenticator("authenticator:6001", iam_role=environ.get("TEST_IAM_ROLE"))
conns = ThreadedConnectionPool(1, 5, "", authenticator=auth)
conn = conns.getconn()
print("Connection Established")


def test():
    conn.cursor().execute("SELECT 1")


test()
