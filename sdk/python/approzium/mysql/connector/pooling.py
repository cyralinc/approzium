from mysql.connector.pooling import MySQLConnectionPool

from ._connect import _parse_kwargs, _patch_MySQLConnection


class MySQLConnectionPool(MySQLConnectionPool):
    def set_config(self, **kwargs):
        kwargs = _parse_kwargs(kwargs)
        super(MySQLConnectionPool, self).set_config(**kwargs)

    def add_connection(self, cnx=None):
        with _patch_MySQLConnection(include_pooling=True):
            super().add_connection(cnx)
