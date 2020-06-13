import psycopg2
from approzium import Authenticator, set_default_authenticator
from approzium.psycopg2 import connect
import logging

logger = logging.getLogger('approzium')
logger.setLevel(logging.DEBUG)


auth = Authenticator('authenticator:1234', 'arn:aws:iam::accountid:role/rolename')
set_default_authenticator(auth)
conn = connect("")  # or connect("", authenticator=auth)
print('Connection Established')

def test():
    conn.cursor().execute('SELECT 1')

test()
