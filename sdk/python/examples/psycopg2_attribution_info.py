"""This example shows usage of the AuthClient.attribution_info feature. It uses
a Psycopg2 connection but the same functionality is available in any supported
database driver.
"""
import approzium
from approzium.psycopg2 import connect

auth = approzium.AuthClient("authenticator:6001")
print(auth.attribution_info)
# {'authenticator_address': 'authenticator:6001',
#  'iam_arn': 'arn:aws:iam::*******:user/****',
# 'authenticated': False,
# 'num_connections': 0
# ... additional info if available (ex: EC2 instance metadata)
# }
approzium.default_auth_client = auth
dsn = "host=dbmd5 dbname=db user=bob"
conn = connect(dsn)
print(auth.attribution_info)
# {'authenticator_address': 'authenticator:6001',
#  'iam_arn': 'arn:aws:iam::*******:user/****',
# 'authenticated': True,
# 'num_connections': 1
# ...
# }
