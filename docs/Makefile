website-dev:
	yarn install
	gatsby clean
	gatsby develop

website-build:
	yarn install
	gatsby clean
	gatsby build
	gatsby serve

# this is meant to run on the Netlify build env
python-docs:
	# ensure sphinx
	pip3 install sphinx
	# other optional dependencies, because autodoc needs to import whole library
	pip3 install psycopg2 asyncpg
	cd ../sdk/python/ && pip3 install -e .
	make -C ../sdk/python/docs/ html
	mkdir -p static/api
	cp -r ../sdk/python/docs/build/html static/api/python
