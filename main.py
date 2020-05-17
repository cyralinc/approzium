#!/usr/bin/env python3
import psycopg2

connstring = "usedbauth=yes host=pc-testing-2.cd6z0yimd7qu.us-west-2.rds.amazonaws.com password=password user=bob dbname=finance";
conn = psycopg2.connect(connstring)
print('Connection established')
