from . import connect
import psycopg2.pool


class AbstractConnectionPool(psycopg2.pool.AbstractConnectionPool):
    def _connect(self, key=None):
        # override connect method
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
