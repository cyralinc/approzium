from os import environ

import approzium
from approzium.psycopg2 import connect

auth = approzium.AuthClient("authenticator:6000", iam_role=environ.get("TEST_IAM_ROLE"))
approzium.default_auth_client = auth
conn = connect("")  # or connect("", authenticator=auth)
print("Connection Established")


def test():
    conn.cursor().execute("SELECT 1")


test()
