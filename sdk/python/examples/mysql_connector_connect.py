from os import environ

from approzium import AuthClient
from approzium.mysql.connector import connect
from approzium.mysql.connector.pooling import MySQLConnectionPool

auth = AuthClient(
    "authenticator:6001",
    trusted_certs=environ.get("TEST_CERT_DIR") + "/approzium.pem",
    client_cert=environ.get("TEST_CERT_DIR") + "/client.pem",
    client_key=environ.get("TEST_CERT_DIR") + "/client.key",
    disable_tls=environ.get("APPROZIUM_DISABLE_TLS"),
)
conn = connect(user="bob", authenticator=auth, host="dbmysql", use_pure=True)
print("Connection Established")

cur = conn.cursor()
cur.execute("SELECT 1")
result = next(cur)
print(result)

cnxpool = MySQLConnectionPool(
    pool_name="mypool", pool_size=3, user="bob", host="dbmysql", authenticator=auth
)
print("Connection Pool Established")
conn = cnxpool.get_connection()
cur = conn.cursor()
cur.execute("SELECT 1")
print(result)
