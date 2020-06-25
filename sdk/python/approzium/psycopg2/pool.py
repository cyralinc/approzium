import psycopg2.pool

from . import connect


class AbstractConnectionPool(psycopg2.pool.AbstractConnectionPool):
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


class SimpleConnectionPool(AbstractConnectionPool, psycopg2.pool.SimpleConnectionPool):
    pass


class ThreadedConnectionPool(
    AbstractConnectionPool, psycopg2.pool.ThreadedConnectionPool
):
    pass
