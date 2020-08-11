import asyncio

import asyncpg
from asyncpg.connection import Connection

from ._asyncpg_connect import connect, new__connect_addr  # approzium's connect method


class _ApproziumPool(asyncpg.pool.Pool):
    async def _get_new_connection(self):
        if self._working_addr is None:
            # First connection attempt on this pool.
            # 1) use our connect method instead of asyncpg's
            con = await connect(
                *self._connect_args, loop=self._loop, **self._connect_kwargs
            )

            self._working_addr = con._addr
            self._working_config = con._config
            self._working_params = con._params
            # 2) after first connection, store class instance because it
            # contains authenticator
            self._connection_class = con.__class__

        else:
            # We've connected before and have a resolved address,
            # and parsed options and config.

            # 3) use our own _connect_addr instead of asyncpg's
            con = await new__connect_addr(
                loop=self._loop,
                addr=self._working_addr,
                timeout=self._working_params.connect_timeout,
                config=self._working_config,
                params=self._working_params,
                connection_class=self._connection_class,
            )

        if self._init is not None:
            try:
                await self._init(con)
            except (Exception, asyncio.CancelledError) as ex:
                # If a user-defined `init` function fails, we don't
                # know if the connection is safe for re-use, hence
                # we close it.  A new connection will be created
                # when `acquire` is called again.
                try:
                    # Use `close()` to close the connection gracefully.
                    # An exception in `init` isn't necessarily caused
                    # by an IO or a protocol error.  close() will
                    # do the necessary cleanup via _release_on_close().
                    await con.close()
                finally:
                    raise ex

        return con


def create_pool(
    dsn=None,
    *,
    min_size=10,
    max_size=10,
    max_queries=50000,
    max_inactive_connection_lifetime=300.0,
    setup=None,
    init=None,
    loop=None,
    authenticator=None,
    **connect_kwargs,
):
    """Create an Asyncpg connection pool through Approzium authentication.
    Takes same arguments as ``asyncpg.create_pool`` in addition to the
    `authenticator` argument

    :return: An instance of :class:`~approzium.asyncpg.pool._ApproziumPool`.

    Example:

    .. code-block:: python

        >>> import approzium
        >>> from approzium.asyncpg import create_pool
        >>> auth = approzium.AuthClient("myauthenticator.com:6001", disable_tls=True)
        >>> pool = await create_pool(user='postgres', authenticator=auth)
        >>> con = await pool.acquire()
        >>> try:
        ...     await con.fetch('SELECT 1')
        ... finally:
        ...     await pool.release(con)
    """
    return _ApproziumPool(
        dsn,
        connection_class=Connection,
        min_size=min_size,
        max_size=max_size,
        max_queries=max_queries,
        loop=loop,
        setup=setup,
        init=init,
        max_inactive_connection_lifetime=max_inactive_connection_lifetime,
        authenticator=authenticator,
        **connect_kwargs,
    )
