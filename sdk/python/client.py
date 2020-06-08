from approzium.psycopg2 import connect


conn = connect("",
        authenticator='authenticator:1234',
        iam_arn='arn:aws:iam::403019568400:role/dev'
)
print('Connection Established')
