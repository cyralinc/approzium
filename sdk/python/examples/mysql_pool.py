from os import environ

import approzium
from approzium.mysql.connector.pooling import MySQLConnectionPool

auth = approzium.AuthClient("authenticator:6000", iam_role=environ.get("TEST_IAM_ROLE"))
dbconfig = {"user": "bob", "authenticator": auth, "host": "dbmysql", "use_pure": True}

cnxpool = MySQLConnectionPool(pool_name="mypool", pool_size=3, **dbconfig)
print("Connection Established")
conn = cnxpool.get_connection()
cur = conn.cursor()
cur.execute("SELECT 1")
result = next(cur)
print(result)
