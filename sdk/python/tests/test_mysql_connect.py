from os import environ

import approzium
import pytest
from approzium.mysql.connector import connect
from approzium.mysql.connector.pooling import MySQLConnectionPool

# use Psycopg2 defined test environment variables
connopts = {
    "user": environ["PSYCOPG2_TESTDB_USER"],
    "host": "dbmysqlsha1",
    "use_pure": True,
}


@pytest.mark.parametrize("auth", pytest.authclients)
def test_connect(auth):
    conn = connect(**connopts, authenticator=auth)
    cur = conn.cursor()
    cur.execute("SELECT 1")
    result = next(cur)
    assert result == (1,)


@pytest.mark.parametrize("auth", pytest.authclients)
def test_pooling(auth):
    approzium.default_auth_client = auth
    cnxpool = MySQLConnectionPool(pool_name="testpool", pool_size=3, **connopts)
    conn = cnxpool.get_connection()
    cur = conn.cursor()
    cur.execute("SELECT 1")
    result = next(cur)
    assert result == (1,)
