deps:
	python3 -m grpc_tools.protoc -I../../authenticator/protos --python_out=approzium/_protos --grpc_python_out=approzium/_protos ../../authenticator/protos/authenticator.proto

download_psycopg2_testsuite:
	$(eval TMP := $(shell mktemp -d))
	git clone https://github.com/psycopg/psycopg2.git $(TMP)
	mkdir -p tests/pg2_testsuite/
	mv -f $(TMP)/tests/* tests/pg2_testsuite/
	rm -rf $(TMP)

# To run the test suite, first set the environment variables:
# - TEST_IAM_ROLE
# - PSYCOPG2_TESTDB
# - PSYCOPG2_TESTDB_HOST
# - PSYCOPG2_TESTDB_PORT
# - PSYCOPG2_TESTDB_USER=bob
# You can automate setting setting up the test database servers and environment
# variables using docker-compose by running `make test-e2e` in the project root
# directory.
test:	download_psycopg2_testsuite
	pip3 install -e .
	cd tests && python3 run_pg2_testsuite.py
	rm -rf tests/pg2_testsuite/
