import psycopg2
from approzium.psycopg2 import connect
from approzium import set_authenticator, set_iam_role
import logging

logger = logging.getLogger('approzium')
logger.setLevel(logging.DEBUG)


set_iam_role('arn:aws:iam::accountid:role/rolename')
set_authenticator('authenticator:1234')
conn = connect("")
print('Connection Established')

def test():
    conn.cursor().execute('SELECT 1')

test()
