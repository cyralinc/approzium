ssl-certs:
	cd ssl && ./gen_cert.sh

test-e2e: ssl-certs
	docker-compose up -d
	docker-compose run psycopg2-testsuite-md5 make test
	docker-compose run psycopg2-testsuite-sha256 make test

test: test-e2e
