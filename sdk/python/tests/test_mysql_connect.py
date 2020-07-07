from os import environ

import approzium
from approzium.mysql.connector import connect
from approzium.mysql.connector.pooling import MySQLConnectionPool

auth = approzium.AuthClient("authenticator:6001", iam_role=environ.get("TEST_IAM_ROLE"))
# use Psycopg2 defined test environment variables
connopts = {
    "user": environ["PSYCOPG2_TESTDB_USER"],
    "host": "dbmysqlsha1",
    "use_pure": True,
}


def test_connect():
    conn = connect(**connopts, authenticator=auth)
    cur = conn.cursor()
    cur.execute("SELECT 1")
    result = next(cur)
    assert result == (1,)


def test_pooling():
    approzium.default_auth_client = auth
    cnxpool = MySQLConnectionPool(pool_name="testpool", pool_size=3, **connopts)
    conn = cnxpool.get_connection()
    cur = conn.cursor()
    cur.execute("SELECT 1")
    result = next(cur)
    assert result == (1,)
