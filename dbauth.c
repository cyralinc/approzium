#include "dbauth.h"

PGconn *DBAconnectdb(const char *conninfo) {
    PGconn *conn = PQconnectdb(conninfo);
    return conn;
} 
