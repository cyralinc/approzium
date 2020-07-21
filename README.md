# approzium

![test](https://github.com/cyralinc/approzium/workflows/test/badge.svg)
![lint](https://github.com/cyralinc/approzium/workflows/lint/badge.svg)
[![Documentation Status](https://readthedocs.org/projects/approzium/badge/?version=latest)](http://approzium.readthedocs.io/?badge=latest)

Approzium provides SDKs that allow you to authenticate to a database without ever having access to its password. Your
identity is provided through the platform on which you're running.

----

**Please note**: We take Approzium's security and our users' trust very seriously. If you believe you have found a security issue in Approzium, _please responsibly disclose_ by contacting us at [security@cyral.com](mailto:security@cyral.com).

See the [SECURITY](.github/SECURITY.md) guide for more details.

----

We currently support AWS for identity, and have a Python SDK for Postgres drivers. This project is under active development, please
do stay tuned for more identity platforms, databases, and SDK languages.

## Docs

See https://approzium.com/ for a Quick Start, or elaboration on the architecture and API.

## Support

For questions, please either open a Github issue, or visit us in our public Slack channel.

To visit us in Slack, use [this invite](https://join.slack.com/t/approzium/shared_invite/zt-fg9bdcfa-H9YFnlg3XeosKyMIYadmcg). 
Then venture to [# help-and-questions](https://app.slack.com/client/T013VTLTTJ5/C013FTJPAN9).
Our developers frequent our Slack forum, but are not in it at all times. Please be patient, we will lend assistance as 
soon as we can!

## Developing

We welcome community contributions!

We use `docker-compose.yml` to quickly and easily provide you with a development environment that mimics real life.
To spin up an end-to-end development environment based in Docker:

- Ensure you have [Docker](https://www.docker.com/) installed with Buildkit support (Docker 18.09 or higher)
- In your local environment, run `$ aws configure` and add an access key and secret.
- Run `$ make dc-build`. This will build the authenticator and development Docker images.
- Run `$ docker-compose up`. This will run the authenticator with a Vault backend and will run test database servers (Postgres and MySQL).
- In another window, `$ make dev`. This will start a shell in the development environment.
- You now have a full development and testing environment!
- For example, to use our Python SDK to create an Approzium connection to a Postgres server:
    * Create an Approzium path in the test Vault backend: `$ make enable-vault-path`
    * Give your AWS-identity access to the test server: `$ make seed-vault-addr ADDR=dbmd5:5432`
    * Create a connection: `$ cd sdk/python/examples && poetry run python3 psycopg2_connect.py`.

### Testing

Our end-to-end tests take a few minutes to run. Please run them once locally before you submit a PR.

To run the end-to-end test, from our home directory:
- Run `make test`. That's it!

## Credits

This project is brought to you by [Cyral](https://www.cyral.com/), who wishes to give back to the Open Source community.
