# approzium
credential-less database authentication for services

## Developing

We welcome community contributions!

We use `docker-compose.yml` to quickly and easily provide you with a development environment that mimics real life. Please check it out for helpful hints on how to reach one container from another. To run the end-to-end test, from our home directory:
- Ensure you have [Docker](https://www.docker.com/) installed with Buildkit support (Docker 18.09 or higher)
- Run `make test`. That's it!

To drop into a Bash shell into the development environment, run `make dev`. This will automatically run everything you need in order to test and debug your code.
- Ensure you have the latest authenticator binary: `$ cd authenticator && GOOS=linux GOARCH=amd64 go build && cd ..`.
- In your local environment, run `$ aws configure` and add an access key and secret.
  - Make sure the access key and secret you configure can assume at least one role.
- In one window, `$ docker-compose up`.
- In another window, `$ make dev`
- Export an environment variable for the role you're testing with: `$ export TEST_IAM_ROLE=export TEST_IAM_ROLE=arn:aws:iam::123456789012:role/AssumableRole`.
- To use our Python SDK to shoot a request at the authenticator, run
  `$ PGHOST=dbmd5 PGUSER=bob PGDATABASE=db python3 sdk/python/client.py`.
