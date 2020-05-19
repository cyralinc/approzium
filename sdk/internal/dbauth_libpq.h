/* this file is not meant for the end user. it is meant to be included by libpq
 * source. It includes functions that communicate with Authenticator service on
 * behalf of libpq. */

#ifndef __DBAUTH_LIBPQ_H
#define __DBAUTH_LIBPQ_H
#ifdef __cplusplus
extern "C" {
#endif
char *dbauth_get_hashed_password(char *md5Salt);
#ifdef __cplusplus
}
#endif
#endif
