# SSL target is required to run Postgres DBs with SSL support
test-e2e:  ssl/rootCA.key
	$(pg2_testsuite_env) docker-compose -f docker-compose.yml -f docker-compose.test.yml up --exit-code-from tests

test: test-e2e

# Starts a bash shell in the `test` service
dev:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run tests bash

# USER INPUT TEST PARAMETERS
TEST_IAM_ROLE=arn:aws:iam::403019568400:role/dev
TEST_DBHOSTS=dbmd5 dbsha256
TEST_DB=db
TEST_DBPORT=5432
TEST_DBPASS=password
TEST_DBUSER=bob

# DERIVE OTHER PARAMETERS
vault_secret = { $\
"password": "$(TEST_DBPASS)", $\
"iam_roles": [ $\
	"$(TEST_IAM_ROLE)" $\
] $\
}
pg2_testsuite_env = TEST_IAM_ROLE=$(TEST_IAM_ROLE) PSYCOPG2_TESTDB=$(TEST_DB) $\
		PSYCOPG2_TESTDB_HOST=$(TEST_DBHOST) PSYCOPG2_TESTDB_PORT=$(TEST_DBPORT)
		PSYCOPG2_TESTDB_USER=$(TEST_DBUSER)


# Generates self-signed certificates that can be used to run Postgres DBs with SSL
ssl/rootCA.key:
	cd ssl && ./gen_cert.sh

# Following targets are called by the `tests` Docker compose service
enable-vault-path:
	vault secrets enable -path=approzium -version=1 kv

run-testsuite: enable-vault-path
	for HOST in $(TEST_DBHOSTS); do \
		vault kv put approzium/$$HOST:$(TEST_DBPORT) $(TEST_DBUSER)='$(vault_secret)'; \
		echo '###### Testing with DBHOST' $$HOST 'SSL=ON #####'; \
		PGSSLMODE=require PSYCOPG2_TESTDB_HOST=$$HOST make -C sdk/python/ test; \
		echo '###### Testing with DBHOST' $$HOST 'SSL=OFF #####'; \
		PGSSLMODE=disable PSYCOPG2_TESTDB_HOST=$$HOST make -C sdk/python/ test; \
	done
