import psycopg2
from approzium.psycopg2 import connect
from approzium import set_authenticator, set_iam_role


set_iam_role('arn:aws:iam::accountid:role/rolename')
set_authenticator('authenticator:1234')
conn = connect("")
print('Connection Established')
