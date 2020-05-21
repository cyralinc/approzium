/* this file is not meant for the end user. it is meant to be included by libpq
 * source. It includes functions that communicate with Authenticator service on
 * behalf of libpq. */

#ifndef __DBAUTH_LIBPQ_H
#define __DBAUTH_LIBPQ_H
#ifdef __cplusplus
extern "C" {
#endif
    /* this function connects to the authenticator service and gets the MD5 hash
     * using the given salt. If the request succeeds, the hash is stored into the
     * result_buffer and a status code of 1 is returned. if the request fails, a 
     * status code of 0 is returned. */
int dbauth_get_hashed_password(char *user, char *md5Salt, int saltlen, char *result_buffer);
#ifdef __cplusplus
}
#endif
#endif
