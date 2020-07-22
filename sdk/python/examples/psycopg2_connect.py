import approzium
from approzium import AuthClient
from approzium.psycopg2 import connect
from approzium.psycopg2.pool import ThreadedConnectionPool

auth = AuthClient(
    "authenticator:6001",
    # This is insecure, see https://approzium.org/configuration for proper use.
    disable_tls=True,
)
dsn = "host=dbmd5 dbname=db user=bob"
conn = connect(dsn, authenticator=auth)
print("Connection Established")

approzium.default_auth_client = auth
conns = ThreadedConnectionPool(1, 5, dsn)
conn = conns.getconn()
print("Connection Pool Established")
