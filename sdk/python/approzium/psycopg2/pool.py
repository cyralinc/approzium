import psycopg2.pool

from . import connect


class _AbstractConnectionPool(psycopg2.pool.AbstractConnectionPool):
    def _connect(self, key=None):
        # originally, this line uses psycopg2.connect. we change it to use our
        # connect method instead
        conn = connect(*self._args, **self._kwargs)
        if key is not None:
            self._used[key] = conn
            self._rused[id(conn)] = key
        else:
            self._pool.append(conn)
        return conn


class SimpleConnectionPool(_AbstractConnectionPool, psycopg2.pool.SimpleConnectionPool):
    """A Psycopg2 pool that can't be shared across different threads.

    Example:

    .. code-block:: python

        >>> import approzium
        >>> from approzium.psycopg2.pool import SimpleConnectionPool
        >>> auth = approzium.AuthClient("myauthenticator.com:6001", disable_tls=True)
        >>> pool = SimpleConnectionPool("host=DB.com dbname=mydb", authenticator=auth)
        >>> conn = poll.getconn()
    """

    pass


class ThreadedConnectionPool(
    _AbstractConnectionPool, psycopg2.pool.ThreadedConnectionPool
):
    """A Psycopg2 pool that works with the threading module"""

    pass
