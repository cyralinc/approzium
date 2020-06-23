from os import environ
import asyncio
import asyncpg

from approzium import Authenticator
from approzium.asyncpg import connect

auth = Authenticator("authenticator:1234", iam_role=environ.get("TEST_IAM_ROLE"))


async def run():
    conn = await connect(user='bob', database='db', host='dbmd5',
                         authenticator=auth)
    print('Connection Established!')
    values = await conn.fetch('''SELECT 1''')
    await conn.close()

loop = asyncio.get_event_loop()
loop.run_until_complete(run())
