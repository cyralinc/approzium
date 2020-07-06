User Guide
**********

Installation
------------

The last stable release is available on PyPI and can be installed with ``pip``::

    $ python3 -m pip install approzium

Requirements
-------------

* Python -- one of the following:

    - CPython_ >= 3.5

.. _CPython: http://www.python.org/

Supported Database Drivers
--------------------------


The following database driver libraries are supported:


      +------------+--------------------+----------------------------------------------------------+-------------------------------------------------------------+
      | Database   | Driver             | Authentication Methods                                   | Notes                                                       |
      +============+====================+==========================================================+=============================================================+
      | Postgres   | Psycopg2_          | MD5 (Postgres default) and SCRAM-SHA-256 authentication  |                                                             |
      +------------+--------------------+----------------------------------------------------------+-------------------------------------------------------------+
      | Postgres   | Asyncpg_           | Same as above                                            |                                                             |
      +------------+--------------------+----------------------------------------------------------+-------------------------------------------------------------+
      | MySQL      | `MySQL Connector`_ | `Native password authentication`_                        | Currently, only the pure Python implementation is supported |
      +------------+--------------------+----------------------------------------------------------+-------------------------------------------------------------+

.. _Psycopg2: https://github.com/psycopg/psycopg2
.. _Asyncpg: https://github.com/MagicStack/asyncpg
.. _MySQL Connector: https://dev.mysql.com/doc/connector-python/en/
.. _Native password authentication: https://dev.mysql.com/doc/refman/8.0/en/native-pluggable-authentication.html


Usage
-----

Approzium Python SDK is designed to have a small footprint on the source code of your application.

1. The first step in creating an Approzium database connection is instantiating an :class:`approzium.AuthClient`:

    .. code-block:: python

        import approzium
        auth = approzium.AuthClient("authenticator_service_host:port")

2. By default, the AuthClient automatically detects the environment that the service is running in. Currently, only AWS-based identity is supported, so it will detect that.

3. Set this auth client to be the default one. This can help you reduce code bloat:

    .. code-block:: python

        approzium.default_auth_client = auth

4. Create a connection! The way you create a connection is extremely similar to existing code. All you have to do is prepend ``approzium.`` to the import path. For example, if you are creating a Psycopg2 connection, instead of ``psycopg2.connect`` you would use ``approzium.psycopg2.connect``. It's that easy!

    .. code-block:: python

        from approzium.psycopg2 import connect
        conn = connect(dbname="test", user="postgres", host="host.com")
        # conn is now a psycopg2.Connection object
