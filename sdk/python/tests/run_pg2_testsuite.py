"""
This script runs the Psycopg2 test suite without a running authenticator
service. Therefore, it only tests the client SDK and a password is needed
"""
import logging
import os
import sys
import unittest
from os.path import abspath
from unittest import SkipTest

import pg2_testsuite  # noqa: E402 F401
import psycopg2

from approzium import Authenticator, set_default_authenticator

# once get hash is mocked, import connect method
from approzium.psycopg2 import connect as approzium_pg2_connect

sys.path.append(abspath("pg2_testsuite"))


try:
    test_iam_role = os.environ["TEST_IAM_ROLE"]
except KeyError:
    print("Please set env var 'TEST_IAM_ROLE' to a valid value")
    sys.exit(1)

auth = Authenticator("authenticator:6001", test_iam_role)
set_default_authenticator(auth)

# replace connect method
psycopg2.connect = approzium_pg2_connect


def skip_tests(ids, pg2_suites):
    filtered_suite = unittest.TestSuite()
    # the way the Psycopg2 test suite is setup, there are many layers of
    # organization
    for suites in pg2_suites:
        for suite in suites:
            for test in suite:
                if test.id() in skipped_tests:

                    def skip_method():
                        raise SkipTest("skipped by Approzium")

                    setattr(test, test._testMethodName, skip_method)
                filtered_suite.addTest(test)
    return filtered_suite


# tests that are accepted to not work with Approzium
skipped_tests = [
    "pg2_testsuite.test_async.AsyncTests.test_async_subclass",
    "pg2_testsuite.test_module.ConnectTestCase.test_no_keywords",
    "pg2_testsuite.test_module.ConnectTestCase.test_dsn",
    "pg2_testsuite.test_module.ConnectTestCase.test_async",
    "pg2_testsuite.test_module.ConnectTestCase.test_factory",
    "pg2_testsuite.test_module.ConnectTestCase.test_flush_on_write",
]
pg2_suite = unittest.TestLoader().loadTestsFromName("pg2_testsuite.test_suite")
suite = skip_tests(skipped_tests, pg2_suite)

logging.getLogger().setLevel(logging.ERROR)
unittest.TextTestRunner(verbosity=1).run(suite)
