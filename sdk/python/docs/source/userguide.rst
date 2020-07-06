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


      +------------+------------+----------------------------------------------------------+
      | Database   | Driver     | Authentication Methods                                   |
      +============+============+==========================================================+
      | Postgres   | Psycopg2   | MD5 (Postgres default) and SCRAM-SHA-256 authentication  |
      +------------+------------+----------------------------------------------------------+
