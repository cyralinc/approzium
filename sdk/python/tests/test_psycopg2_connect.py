from os import environ
from select import select

import psycopg2
import pytest
import approzium
from approzium.psycopg2 import connect
from approzium.psycopg2.pool import SimpleConnectionPool, ThreadedConnectionPool

auth = approzium.AuthClient("authenticator:1234", iam_role=environ.get("TEST_IAM_ROLE"))
# use Psycopg2 defined test environment variables
connopts = {
    "user": environ["PSYCOPG2_TESTDB_USER"],
    "dbname": environ["PSYCOPG2_TESTDB"],
    "port": environ["PSYCOPG2_TESTDB_PORT"],
}


# waits for an async connection to start or finish execution
# source: https://www.psycopg.org/docs/advanced.html
def wait(conn):
    while True:
        state = conn.poll()
        if state == psycopg2.extensions.POLL_OK:
            break
        elif state == psycopg2.extensions.POLL_WRITE:
            select([], [conn.fileno()], [])
        elif state == psycopg2.extensions.POLL_READ:
            select([conn.fileno()], [], [])
        else:
            raise psycopg2.OperationalError("poll() returned %s" % state)


@pytest.mark.parametrize("dbhost", ["dbmd5", "dbsha256"])
@pytest.mark.parametrize("sslmode", ["require", "disable"])
@pytest.mark.parametrize("async_", [1, 0])
def test_connect(dbhost, sslmode, async_):
    conn = connect(
        **connopts, host=dbhost, sslmode=sslmode, async_=async_, authenticator=auth
    )
    if async_:
        wait(conn)
    cur = conn.cursor()
    cur.execute("SELECT 1")
    if async_:
        wait(conn)
    assert cur.fetchone() == (1,)


@pytest.mark.parametrize("dbhost", ["dbmd5", "dbsha256"])
@pytest.mark.parametrize("sslmode", ["require", "disable"])
@pytest.mark.parametrize("async_", [1, 0])
@pytest.mark.parametrize("Pool", [ThreadedConnectionPool, SimpleConnectionPool])
def test_pool(dbhost, sslmode, async_, Pool):
    approzium.default_auth_client = auth
    conns = Pool(
        1,
        5,
        "",
        **connopts,
        host=dbhost,
        sslmode=sslmode,
        async_=async_,
    )
    conn = conns.getconn()
    if async_:
        wait(conn)
    cur = conn.cursor()
    cur.execute("SELECT 1")
    if async_:
        wait(conn)
    assert cur.fetchone() == (1,)
