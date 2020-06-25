from os import environ

import approzium
import pytest
from approzium.asyncpg import connect
from approzium.asyncpg.pool import create_pool

auth = approzium.AuthClient("authenticator:1234", iam_role=environ.get("TEST_IAM_ROLE"))
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


@pytest.mark.parametrize("host", ["dbmd5", "dbsha256"])
@pytest.mark.parametrize("sslmode", ["disable", "require"])
@pytest.mark.asyncio
async def test_pool(host, sslmode):
    approzium.default_auth_client = auth
    # set SSL mode using env variable because there is no better way
    environ["PGSSLMODE"] = sslmode
    pool = await create_pool(**connopts, host=host)
    async with pool.acquire() as conn:
        await conn.fetch("""SELECT 1""")
