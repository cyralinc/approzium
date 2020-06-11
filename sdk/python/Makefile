deps:
	python3 -m grpc_tools.protoc -I../../authenticator/protos --python_out=approzium/protos --grpc_python_out=approzium/protos ../../authenticator/protos/authenticator.proto

psycopg2_testsuite:
	$(eval TMP := $(shell mktemp -d))
	git clone https://github.com/psycopg/psycopg2.git $(TMP)
	mkdir -p tests/pg2_testsuite/
	mv -f $(TMP)/tests/* tests/pg2_testsuite/
	rm -rf $(TMP)

test:	psycopg2_testsuite
	pip3 install -e .
	cd tests && python3 run_pg2_testsuite.py
	rm -rf tests/pg2_testsuite/
