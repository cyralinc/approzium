# approzium

![test](https://github.com/cyralinc/approzium/workflows/test/badge.svg)
![lint](https://github.com/cyralinc/approzium/workflows/lint/badge.svg)
[![Documentation Status](https://readthedocs.org/projects/approzium/badge/?version=latest)](http://approzium.readthedocs.io/?badge=latest)

Approzium is a tool that provides:
- Password-less database authentication
- Authentication through your cloud-provider's built-in identity
- Highly security-oriented logging and metrics

Its aim is to prevent data breaches, and to help you detect them promptly if they do occur or are attempted.

----

**Please note**: We take Approzium's security and our user's trust very seriously. If you believe you have found a security issue in Approzium, _please responsibly disclose_ by contacting us at [security@cyral.com](mailto:security@cyral.com).

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
- In your local environment, run `$ aws configure` and add an access key and a secret. Also,
make sure that you have the `AWS_REGION` environment variable set, informing the AWS region that will be used. For instance:
```
export AWS_REGION=us-east-1
```
- Then run `$ make dev-env`. This will build the authenticator and development Docker images. Also, it will run the authenticator with a Vault backend and the test database servers (Postgres and MySQL).
- In another window, `$ make dev`. This will start a shell in the development environment.
- You now have a full development and testing environment!
- For example, to use our Python SDK to create an Approzium connection to a Postgres server:
    * Create an Approzium path in the test Vault backend: `$ make enable-vault-path`
    * Give your AWS-identity access to the test server: `$ make seed-vault-addr ADDR=dbmd5:5432`
    * Create a connection: `$ cd sdk/python/examples && poetry run python3 psycopg2_connect.py`.

### Testing

Our end-to-end tests take a few minutes to run. Please run them once locally before you submit a PR.

To run the tests, first you will need to:
- Create an AWS `Role` (E.g. ApproziumTestAssumableRole) thats going to be used during the tests.
- Ensure that you are using an AWS `User` with at least the following permissions:
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "secretsmanager:CreateSecret",
                "secretsmanager:GetSecretValue",
                "secretsmanager:UpdateSecret",
                "secretsmanager:DeleteSecret",
                "secretsmanager:PutSecretValue"
            ],
            "Resource": "arn:aws:secretsmanager:us-east-2:<some-account-id>:secret:approzium/*"
        },
        {
            "Effect": "Allow",
            "Action": "sts:AssumeRole",
            "Resource": "arn:aws:iam::<some-account-id>:role/ApproziumTestAssumableRole"
        }
    ]
}
```
- Set the `AWS_REGION` and the `TEST_ASSUMABLE_ARN` environment variables, for instance:
```
export AWS_REGION=us-east-1 && \
export TEST_ASSUMABLE_ARN=arn:aws:iam::<some-account-id>:role/ApproziumTestAssumableRole
```
Then, to run the end-to-end tests, from our home directory:
- Run `make test`. That's it!

## Credits

This project is brought to you by [Cyral](https://www.cyral.com/), who wishes to give back to the Open Source community.
