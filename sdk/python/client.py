from approzium.psycopg2 import connect


conn = connect('authenticator=authenticator:1234')
print('Connection Established')
