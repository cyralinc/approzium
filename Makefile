build-env:
	docker build -t dbauth-dev .
run-env:
	docker run -it -v "$(PWD)":/usr/src/dbauth --rm --name dbauth-dev dbauth-dev bash