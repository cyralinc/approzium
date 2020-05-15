#include <stdio.h>
#include <stdlib.h>
#include "dbauth.h"

int main() {
    printf("hello\n");
    PGconn *conn = DBAconnectdb(""); // dummy struct
    printf("%d\n", conn->auth_req_received); // can query struct fields
}
