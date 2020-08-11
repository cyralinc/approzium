from os import environ

import pytest

from approzium.pymysql import connect

# use Psycopg2 defined test environment variables
connopts = {
    "user": environ["PSYCOPG2_TESTDB_USER"],
    "host": "dbmysqlsha1",
}


@pytest.mark.parametrize("auth", pytest.authclients)
def test_connect(auth):
    conn = connect(**connopts, authenticator=auth)
    cur = conn.cursor()
    cur.execute("SELECT 1")
    result = cur.fetchone()
    assert result == (1,)
