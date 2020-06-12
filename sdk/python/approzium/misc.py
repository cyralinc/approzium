import threading
import struct
import select
import socket
from contextlib import contextmanager
import os
import queue
import logging


logger = logging.getLogger(__name__)

def read_int32_from_bytes(bytes, index):
    num = struct.unpack("!i", bytes[index : index + 4])[0]
    return num


@contextmanager
def redirect_socket_nowhere(fileno, feedit=None):
    void_socket = new_mim_socket(feedit)
    orig_dest = os.dup(fileno)
    os.dup2(void_socket.fileno(), fileno)
    logger.debug('redirecting socket')
    yield
    os.dup2(orig_dest, fileno)
    logger.debug('restoring original socket')


def mim_conn_listen(clientsocket, feedit):
    child_pid = os.fork()
    if child_pid != 0:
        return
    logger.debug('listening on connection')
    while True:
        rlist, wlist = select.select([clientsocket], [clientsocket], [])[0:2]
        if clientsocket in rlist:
            buf = clientsocket.recv(4096)
            if not buf:
                break
            logger.debug(f'got {buf} and ignoring it')
        if clientsocket in wlist and not feedit is None:
            logging.debug(f'feeding socket client {feedit}')
            clientsocket.sendall(feedit)
            feedit = None

def mim_server_listen(new_server_socket, feedit):
    (clientsocket, address) = new_server_socket.accept()
    thread = threading.Thread(target=mim_conn_listen, args=(clientsocket, feedit), daemon=True)
    thread.start()
    
def new_mim_socket(feedit):
    # picking port as 0 allows OS to pick an available port
    mim_addr = ('', 0)
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as new_server_socket:
        new_server_socket.bind(mim_addr)
        new_server_socket.listen()
        mim_addr = new_server_socket.getsockname()  # to get real port num
        q = queue.Queue()
        thread = threading.Thread(target=mim_server_listen, args=(new_server_socket, feedit), daemon=True)
        thread.start()
        new_conn_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        new_conn_socket.connect(mim_addr)
        # wait until connection is established
        thread.join()
    return new_conn_socket
