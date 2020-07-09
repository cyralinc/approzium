import asyncio
from os import environ

from approzium import AuthClient
from approzium.asyncpg import connect
from approzium.asyncpg.pool import create_pool

auth = AuthClient(
    "authenticator:6001",
    trusted_certs=environ.get("TEST_CERT_DIR") + "/approzium.pem",
    client_cert=environ.get("TEST_CERT_DIR") + "/client.pem",
    client_key=environ.get("TEST_CERT_DIR") + "/client.key",
    disable_tls=environ.get("APPROZIUM_DISABLE_TLS"),
)


async def run():
    conn = await connect(user="bob", database="db", host="host", authenticator=auth)
    print("Connection Established!")
    await conn.fetch("""SELECT 1""")
    await conn.close()

    pool = await create_pool(user="bob", database="db", host="host", authenticator=auth)
    print("Connection Established!")
    async with pool.acquire() as conn:
        await conn.fetch("""SELECT 1""")


loop = asyncio.get_event_loop()
loop.run_until_complete(run())
