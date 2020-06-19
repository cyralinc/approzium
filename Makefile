# Targets that can be run from host machine
# Starts a bash shell in the dev environment
dev:
	$(docker_env) docker-compose $(dc_files) run tests bash
dev-env: dc-build
	$(docker_env) docker-compose up
dc-build: ssl/rootCA.key
	$(docker_env) docker-compose -f docker-compose.yml -f docker-compose.test.yml build
# Runs all tests, including E2E tests
test: run-tests-in-docker

# PARAMETERS USED FOR TESTS
TEST_IAM_ROLE=arn:aws:iam::403019568400:role/dev
TEST_DBHOSTS=dbmd5 dbsha256
TEST_DB=db
TEST_DBPORT=5432
TEST_DBPASS=password
TEST_DBUSER=bob


### Anything below here is implementation details ###

dc_files=-f docker-compose.yml -f docker-compose.test.yml 
# Enable Buildkit in docker commands
docker_env=COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1

run-tests-in-docker:  dc-build  # need SSL certs for Postgres services
	$(docker_env) $(pg2_testsuite_env) docker-compose $(dc_files) up --exit-code-from tests


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
	vault secrets enable -path=approzium -version=1 kv | true
seed-vault-host:  # call this with "make seed-vault-host HOST=foo"
	echo '{"$(TEST_DBUSER)": $(vault_secret)}' | \
		vault write approzium/$(HOST):$(TEST_DBPORT) -

run-testsuite: enable-vault-path run-gotests run-pg2tests

run-gotests:
	echo '###### Running Go tests ######'
	cd authenticator && CGO_ENABLED=1 go test -v -race ./...

run-pg2tests:
	echo '###### Running Psycopg2 test suite ######'
	for HOST in $(TEST_DBHOSTS); do \
		make seed-vault-host HOST=$$HOST \
		echo '###### Testing with DBHOST' $$HOST 'SSL=ON #####'; \
		PGSSLMODE=require PSYCOPG2_TESTDB_HOST=$$HOST make -C sdk/python/ test; \
		echo '###### Testing with DBHOST' $$HOST 'SSL=OFF #####'; \
		PGSSLMODE=disable PSYCOPG2_TESTDB_HOST=$$HOST make -C sdk/python/ test; \
	done
