from os import environ
import approzium
import pytest


def pytest_configure():
    determine_authclients()


def determine_authclients():
    """If TEST_ASSUMABLE_ARN is set, adds an additional Approzium AuthClient
    that uses it"""
    pytest.authclients = []
    base_aws_auth = approzium.AuthClient(
        "authenticator:6001"
    )
    pytest.authclients.append(base_aws_auth)
    if environ.get('TEST_ASSUMABLE_ARN'):
        role_aws_auth = approzium.AuthClient(
            "authenticator:6001"
        )
        pytest.authclients.append(role_aws_auth)
    else:
        print("""Skipping testing using assumable AWS roles because
TEST_ASSUMABLE_ARN is not set""")
