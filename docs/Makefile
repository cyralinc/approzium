website-dev:
	yarn install
	gatsby clean
	gatsby develop

website-build:
	yarn install
	gatsby clean
	gatsby build
	gatsby serve

# this target is meant to run only in the Netlify docs environment
# since we cannot use `sudo` there, we cannot `sudo apt-get install libpq-dev`,
# and thus we cannot install `psycopg2`. Thus, we remove it from the package dependencies 
# before we attempt to generate the documentation.
# the rest of the target installs and runs sphinx, and copies its output to the Gatsby static folder
python-docs:
	pip3 install poetry
	(cd ../sdk/python/ && sed -i '/psycopg2/d' pyproject.toml) || true
	(cd ../sdk/python/ && rm poetry.lock ) || true
	(cd ../sdk/python/ && poetry install) || true
	(cd ../sdk/python/ && poetry add psycopg2-binary) || true
	(cd ../sdk/python/ && poetry add sphinx) || true
	cd ../sdk/python/ && poetry run make -C docs/ html
	mkdir -p static/api
	cp -r ../sdk/python/docs/build/html static/api/python
