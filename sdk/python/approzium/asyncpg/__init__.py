import asyncpg
from asyncpg.connect_utils import _connect_addr
from asyncpg.protocol import Protocol
import logging
from .._postgres import PGAuthClient, parse_msg, construct_msg


logger = logging.getLogger(__name__)


def new_connection_made(authenticator):
    original = Protocol.connection_made
    def connection_made(self, transport):
        transport.authenticator = authenticator
        self.approzium_context = {'authenticator': authenticator,
                                  'transport': transport,
                                  'authclient': None,
                                  'host': authenticator.asyncpgaddress[0],
                                  'port': str(authenticator.asyncpgaddress[1]),
                                  }
        return original(self, transport)
    return connection_made

def hack(authenticator):
    Protocol.connection_made = new_connection_made(authenticator)


original_data_received = Protocol.data_received
def new_data_received(self, data):
    if hasattr(self, 'approzium_context'):
        # this is an approzium transport
        context = self.approzium_context
        # data could contain multiple message, in which case, we want to consume
        # only the first message, while asyncpg consumes the rest
        msg_type, msg_content = parse_msg(data)
        first_msg_bytes = construct_msg(msg_type, msg_content)
        rest_bytes = data[len(first_msg_bytes):]
        read_bytes = lambda: first_msg_bytes  # noqa: E731
        if context['authclient'] is None:
            context['authclient'] = PGAuthClient(
                read_bytes,
                lambda msg: context['transport'].write(msg),
                context['authenticator'],
                context['host'],
                context['port'],
                self.user
            )
        else:
            context['authclient'].read_bytes = read_bytes
        if not context['authclient'].done:
            next(context['authclient'])
            if len(rest_bytes):
                original_data_received(self, rest_bytes)
        else:
            original_data_received(self, data)
Protocol.data_received = new_data_received

original__connect_addr = _connect_addr
def new__connect_addr(*args, **kwargs):
    if kwargs['connection_class'].__name__ == 'ApproziumConnection':
        # store host and port in authenticator instance
        kwargs['connection_class'].authenticator.asyncpgaddress = kwargs['addr']
    return original__connect_addr(*args, **kwargs)
asyncpg.connect_utils._connect_addr = new__connect_addr


def construct_approzium_connection(authenticator):
    class ApproziumConnection(asyncpg.connection.Connection):
        pass
    ApproziumConnection.authenticator = authenticator
    return ApproziumConnection


async def connect(*args, authenticator=None, **kwargs):
    hack(authenticator)
    conn = await asyncpg.connect(*args, **kwargs,
                                 connection_class=construct_approzium_connection(authenticator))
    return conn
