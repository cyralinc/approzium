import approzium
from approzium.psycopg2 import connect
from approzium.psycopg2.pool import ThreadedConnectionPool

auth = approzium.AuthClient("authenticator:6000")
approzium.default_auth_client = auth
dsn = "host=myhody dbname=db user=bob"
conn = connect(dsn)
print("Connection Established")

conns = ThreadedConnectionPool(1, 5, dsn)
conn = conns.getconn()
print("Connection Pool Established")
