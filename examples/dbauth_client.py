#!/usr/bin/env python3
import psycopg2

conn = psycopg2.connect(host='db', authenticator='authenticator:1234', identity='diotim', dbname='db')
print('Connection established')
