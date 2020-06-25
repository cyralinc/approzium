import asyncio
from os import environ

from approzium import Authenticator
from approzium.asyncpg.pool import create_pool

auth = Authenticator("authenticator:1234", iam_role=environ.get("TEST_IAM_ROLE"))


async def run():
    pool = await create_pool(
        user="bob", database="db", host="dbmd5", authenticator=auth
    )
    print("Connection Established!")
    async with pool.acquire() as conn:
        await conn.fetch("""SELECT 1""")


loop = asyncio.get_event_loop()
loop.run_until_complete(run())
