"""
This script runs the Psycopg2 test suite without a running authenticator
service. Therefore, it only tests the client SDK and a password is needed
"""
import unittest
from unittest import SkipTest
import os
from os.path import abspath
import logging
import sys

sys.path.append(abspath("pg2_testsuite"))
import hashlib
import approzium.authenticator


def mock_get_hash(dbhost, dbuser, salt, authenticator):
    password = os.environ["PSYCOPG2_TESTDB_PASSWORD"]
    first_hash = hashlib.md5((password + dbuser).encode("ascii")).hexdigest()
    second_hash = hashlib.md5(first_hash.encode("ascii") + salt).hexdigest()
    return second_hash


approzium.authenticator.get_hash = mock_get_hash
# once get hash is mocked, import connect method
from approzium.psycopg2 import connect
import pg2_testsuite
import psycopg2


# replace connect method
psycopg2.connect = connect


def skip_tests(ids, pg2_suites):
    filtered_suite = unittest.TestSuite()
    # the way the Psycopg2 test suite is setup, there are many layers of organization
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
