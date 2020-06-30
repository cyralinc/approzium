from os import environ

import approzium
from approzium.mysql.connector import connect

auth = approzium.AuthClient("authenticator:6000", iam_role=environ.get("TEST_IAM_ROLE"))
approzium.default_auth_client = auth
conn = connect(user="bob", authenticator=auth, host="dbmysql", use_pure=True)
print("Connection Established")
cur = conn.cursor()
cur.execute("SELECT 1")
result = next(cur)
print(result)
