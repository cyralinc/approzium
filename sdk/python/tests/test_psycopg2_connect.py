import pytest
from os import environ
from approzium import Authenticator
from approzium.psycopg2 import connect


auth = Authenticator("authenticator:1234", iam_role=environ.get("TEST_IAM_ROLE"))
# use Psycopg2 defined test environment variables
connopts = {'user': environ['PSYCOPG2_TESTDB_USER'],
            'dbname': environ['PSYCOPG2_TESTDB'],
            'port': environ['PSYCOPG2_TESTDB_PORT']}

@pytest.mark.parametrize("sslmode", ['require', 'disable'])
@pytest.mark.parametrize("dbhost", ['dbmd5', 'dbsha256'])
def test_connect(dbhost, sslmode):
    conn = connect(**connopts, host=dbhost, sslmode=sslmode, authenticator=auth)
    cur = conn.cursor()
    cur.execute('SELECT 1')
    assert cur.fetchone() == (1, )
