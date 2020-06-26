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
	pip3 install poetry
	(cd ../sdk/python/ && poetry install) || true
	cd ../sdk/python/ && poetry add psycopg2-binary sphinx
	make -C ../sdk/python/docs/ html
	mkdir -p static/api
	cp -r ../sdk/python/docs/build/html static/api/python
