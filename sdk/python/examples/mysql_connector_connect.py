from approzium import AuthClient
from approzium.mysql.connector import connect
from approzium.mysql.connector.pooling import MySQLConnectionPool

auth = AuthClient(
    "authenticator:6000",
    # This is insecure, see https://approzium.org/configuration for proper use.
    disable_tls=True,
)
conn = connect(user="bob", authenticator=auth, host="dbmysql", use_pure=True)
print("Connection Established")

cur = conn.cursor()
cur.execute("SELECT 1")
result = next(cur)
print(result)

cnxpool = MySQLConnectionPool(
    pool_name="mypool",
    pool_size=3,
    user="bob",
    host="dbmysqlsha1",
    authenticator=auth,
    use_pure=True,
)
print("Connection Pool Established")
conn = cnxpool.get_connection()
cur = conn.cursor()
cur.execute("SELECT 1")
print(result)
