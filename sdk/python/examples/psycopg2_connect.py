from os import environ

from approzium import AuthClient
from approzium.psycopg2 import connect
from approzium.psycopg2.pool import ThreadedConnectionPool

auth = AuthClient(
    "authenticator:6001",
    trusted_certs=environ.get("TEST_CERT_DIR") + "/approzium.pem",
    client_cert=environ.get("TEST_CERT_DIR") + "/client.pem",
    client_key=environ.get("TEST_CERT_DIR") + "/client.key",
    disable_tls=environ.get("APPROZIUM_DISABLE_TLS"),
)
dsn = "host=dbmd5 dbname=db user=bob"
conn = connect(dsn, authenticator=auth)
print("Connection Established")

conns = ThreadedConnectionPool(1, 5, dsn)
conn = conns.getconn()
print("Connection Pool Established")
