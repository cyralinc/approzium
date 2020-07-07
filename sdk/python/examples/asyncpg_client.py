import asyncio
from os import environ

from approzium import Authenticator
from approzium.asyncpg import connect

auth = Authenticator("authenticator:6001", iam_role=environ.get("TEST_IAM_ROLE"))


async def run():
    conn = await connect(user="bob", database="db", host="dbmd5", authenticator=auth)
    print("Connection Established!")
    await conn.fetch("""SELECT 1""")
    await conn.close()


loop = asyncio.get_event_loop()
loop.run_until_complete(run())
