import threading
import psycopg2
from psycopg2 import connect


conn = connect("", async=1)

def wait(conn):
        while True:
    state = conn.poll()
    if state == psycopg2.extensions.POLL_OK:
        break
    elif state == psycopg2.extensions.POLL_WRITE:
        select.select([], [conn.fileno()], [])
    elif state == psycopg2.extensions.POLL_READ:
        select.select([conn.fileno()], [], [])
    else:
        raise psycopg2.OperationalError("poll() returned %s" % state)

wait(conn)
cur = conn.cursor()
cur.execute("select 1")
print(conn.info.transaction_status)
