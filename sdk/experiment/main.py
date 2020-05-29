#!/usr/bin/env python3
from ctypes import *
import ctypes.util
import os.path as osp
import psycopg2 as pg2


libpq = CDLL("libpq.so", mode=RTLD_GLOBAL)
libhook = CDLL(osp.abspath("libhook.so"), mode=RTLD_GLOBAL)

conn = pg2.connect("")

libhook.hook2()

@CFUNCTYPE(c_char_p, c_char_p)
def md5_callback(salt):
    print('main.py(md5_callback): salt is', salt.hex())
    hash = "".encode()
    return hash

libhook.register_python_callback(md5_callback)
conn = pg2.connect("")
