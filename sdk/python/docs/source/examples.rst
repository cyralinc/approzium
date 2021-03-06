If any example is broken, or if you’d like to add an example to this page, feel free to raise an issue on our Github repository.

Postgres Examples
-----------------

Example of creating a Psycopg2 single connection and a connection pool

:file:`psycopg2_connect.py`:

.. literalinclude:: ../../examples/psycopg2_connect.py
    :language: python

Example of creating a Asyncpg single connection and a connection pool

:file:`asyncpg_connect.py`:

.. literalinclude:: ../../examples/asyncpg_connect.py
    :language: python


MySQL Examples
--------------

Example of creating a MySQL Connector single connection and a connection pool

:file:`mysql_connector_connect.py`:

.. literalinclude:: ../../examples/mysql_connector_connect.py
    :language: python

Example of creating a PyMySQL single connection:

:file:`pymysql_connect.py`:

.. literalinclude:: ../../examples/pymysql_connect.py
    :language: python

Opentelemetry Integration Examples
----------------------------------

:file:`psycopg2_opentelemetry.py`:

.. literalinclude:: ../../examples/psycopg2_opentelemetry.py
    :language: python

If you are not using Opentelemetry, you can obtain the same attribution info manually:

:file:`psycopg2_attribution_info.py`:

.. literalinclude:: ../../examples/psycopg2_attribution_info.py
