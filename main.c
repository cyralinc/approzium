#include <stdio.h>
#include <stdlib.h>
#include "dbauth.h"

static void exit_nicely(PGconn *conn) {
    PQfinish(conn);
    exit(1);
}

int main() {
    char *connstring = "usedbauth=no host=pc-testing-2.cd6z0yimd7qu.us-west-2.rds.amazonaws.com password=password user=bob dbname=finance";
    PGconn *conn = PQconnectdb(connstring);
    if (PQstatus(conn) != CONNECTION_OK) {
        fprintf(stderr, "Connection to database failed: %s\n", PQerrorMessage(conn));
        exit_nicely(conn);
    }
    fprintf(stderr, "Connection established\n");
    PQfinish(conn);
    return 0;
}
