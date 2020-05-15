#include "dbauth.h"


PGconn *DBAconnectdb(const char *conninfo) {
    return PQconnectdb(conninfo);
}
