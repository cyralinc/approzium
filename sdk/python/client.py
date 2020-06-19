from os import environ
from approzium import AuthClient
from approzium.psycopg2 import connect
import logging
from pprint import pprint

auth = AuthClient("authenticator:1234", iam_role=environ.get("TEST_IAM_ROLE"))
conn = connect("", auth=auth)
print("Connection Established")
# TODO
# add add_attribution_info(key, value)
# capture env vars with certain prefix, names
# try integration with some tracing libs
pprint(auth.attribution_info)
