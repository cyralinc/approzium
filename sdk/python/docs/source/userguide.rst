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


      +------------+--------------------+----------------------------------------------------------+
      | Database   | Driver             | Authentication Methods                                   |
      +============+====================+==========================================================+
      | Postgres   | Psycopg2_          | MD5 (Postgres default) and SCRAM-SHA-256 authentication  |
      +------------+--------------------+----------------------------------------------------------+
      | Postgres   | Asyncpg_           | Same as above                                            |
      +------------+--------------------+----------------------------------------------------------+
      | MySQL      | `MySQL Connector`_ | `Native password authentication`_                        |
      +------------+--------------------+----------------------------------------------------------+

.. _Psycopg2: https://github.com/psycopg/psycopg2
.. _Asyncpg: https://github.com/MagicStack/asyncpg
.. _MySQL Connector: https://dev.mysql.com/doc/connector-python/en/
.. _Native password authentication: https://dev.mysql.com/doc/refman/8.0/en/native-pluggable-authentication.html
