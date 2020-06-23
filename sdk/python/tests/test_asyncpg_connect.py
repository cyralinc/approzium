from os import environ

import pytest
from approzium import Authenticator
from approzium.asyncpg import connect

auth = Authenticator("authenticator:1234", iam_role=environ.get("TEST_IAM_ROLE"))
# use Psycopg2 defined test environment variables
connopts = {
    "user": environ["PSYCOPG2_TESTDB_USER"],
    "database": environ["PSYCOPG2_TESTDB"],
    "port": int(environ["PSYCOPG2_TESTDB_PORT"]),
}


@pytest.mark.parametrize("host", ["dbmd5", "dbsha256"])
@pytest.mark.parametrize("sslmode", ["disable", "require"])
@pytest.mark.asyncio
async def test_connect(host, sslmode):
    # set SSL mode using env variable because there is no better way
    environ["PGSSLMODE"] = sslmode
    conn = await connect(**connopts, host=host, authenticator=auth)
    await conn.fetch("SELECT 1")
