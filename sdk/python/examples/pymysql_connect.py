from approzium import AuthClient
from approzium.pymysql import connect

auth = AuthClient(
    "authenticator:6001",
    # This is insecure, see https://approzium.org/configuration for proper use.
    disable_tls=True,
)
conn = connect(host="dbmysqlsha1", user="bob", db="db", authenticator=auth)
with conn.cursor() as cursor:
    cursor.execute("SELECT 1")
    result = cursor.fetchone()
    print(result)
conn.close()
