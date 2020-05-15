#ifndef DBAUTH_H
#define DBAUTH_H
#include "postgres_fe.h"
#include "libpq-fe.h"
#include "libpq-int.h"

extern PGconn *DBAconnectdb(const char *conninfo);

#endif							/* DBAUTH_H */
