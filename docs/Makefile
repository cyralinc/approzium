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
	source /opt/buildhome/python3.7/bin/activate
	# ensure sphinx
	pip3 install sphinx
	make -C ../sdk/python/docs/ html
	mkdir -p static/api
	cp -r ../sdk/python/docs/build/html static/api/python
