# approzium
Approzium provides SDKs that allow you to authenticate to a database without ever having access to its password. Your
identity is provided through the platform on which you're running.

----

**Please note**: We take Approzium's security and our users' trust very seriously. If you believe you have found a security 
issue in Approzium, _please responsibly disclose_ by contacting us at [security@cyral.com](mailto:security@cyral.com).

----

We currently support AWS for identity, and have a Python SDK for Postgres drivers. This project is under active development, please
do stay tuned for more identity platforms, databases, and SDK languages.

## Docs

See https://approzium.org/ for a Quick Start, or elaboration on the architecture and API.

## Support

We currently answer community questions via a public Slack channel, [# help-approzium](https://app.slack.com/client/T013VTLTTJ5/C0166QPU6LW/thread/C013VTN50JV-1593627063.014500).
Our developers frequent this forum, but are not in it at all times. Please be patient, we will lend assistance as soon
as we can!

## Developing

We welcome community contributions!

We use `docker-compose.yml` to quickly and easily provide you with a development environment that mimics real life. Please check it out for helpful hints on how to reach one container from another. To run the end-to-end test, from our home directory:
- Ensure you have [Docker](https://www.docker.com/) installed with Buildkit support (Docker 18.09 or higher)
- Run `make test`. That's it!

To drop into a Bash shell into the development environment, run `make dev`. This will automatically run everything you need in order to test and debug your code.
- Ensure you have the latest authenticator binary: `$ cd authenticator && GOOS=linux GOARCH=amd64 go build && cd ..`.
- In your local environment, run `$ aws configure` and add an access key and secret.
  - Make sure the access key and secret you configure can assume at least one role.
- Run `$ make dc-build`.
- In one window, `$ docker-compose up`.
- In another window, `$ make dev`.
- Export an environment variable for the role you're testing with: `$ export TEST_IAM_ROLE=arn:aws:iam::123456789012:role/AssumableRole`.
- To use our Python SDK to shoot a request at the authenticator, run
  `$ PGHOST=dbmd5 PGUSER=bob PGDATABASE=db python3 sdk/python/examples/pg2_client.py`.

## Credits

This project is brought to you by [Cyral](https://www.cyral.com/), who wishes to give back to the Open Source community.
