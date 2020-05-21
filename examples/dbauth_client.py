#!/usr/bin/env python3
import psycopg2

conn = psycopg2.connect(host='db', user='bob', usedbauth='yes', dbname='db')
print('Connection established')
