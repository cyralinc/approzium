from os import environ

import approzium
import pytest
from approzium.mysql.connector import connect

auth = approzium.AuthClient("authenticator:6000", iam_role=environ.get("TEST_IAM_ROLE"))
# use Psycopg2 defined test environment variables
connopts = {
    "user": environ["PSYCOPG2_TESTDB_USER"],
    "host": "dbmysql"
}


def test_connect(host, sslmode):
    conn = connect(**connopts, authenticator=auth)
    cur = conn.cursor()
    cur.execute("SELECT 1")
    result = next(cur)
    assert result == (1, )
