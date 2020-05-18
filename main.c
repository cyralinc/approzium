#include <stdio.h>
#include <stdlib.h>
#include "dbauth.h"

static void exit_nicely(PGconn *conn) {
    PQfinish(conn);
    exit(1);
}

int main() {
    ***REMOVED***
    PGconn *conn = PQconnectdb(connstring);
    if (PQstatus(conn) != CONNECTION_OK) {
        fprintf(stderr, "Connection to database failed: %s\n", PQerrorMessage(conn));
        exit_nicely(conn);
    }
    fprintf(stderr, "Connection established\n");
    PQfinish(conn);
    return 0;
}
