env:
	export COMPOSE_DOCKER_CLI_BUILD=1
	export DOCKER_BUILDKIT=1

ssl-certs:
	cd ssl && ./gen_cert.sh

test: ssl-certs env
	# e2e test
	docker-compose up -d
	docker-compose run psycopg2-testsuite-md5 make test
	docker-compose run psycopg2-testsuite-sha256 make test
