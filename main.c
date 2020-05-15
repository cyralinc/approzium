#include "postgres_fe.h"
#include <stdio.h>
#include <stdlib.h>
#include "libpq-fe.h"
#include "libpq-int.h"

int main() {
    printf("hello\n");
    PGconn *conn = PQconnectStart(""); // dummy struct
    printf("%d\n", conn->auth_req_received); // can query struct fields
}
