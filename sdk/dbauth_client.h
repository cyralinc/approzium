/* this file is not meant for the end user. it is meant to be included by libpq
 * source. It includes functions that communicate with Authenticator service on
 * behalf of libpq. */

#ifndef __DBAUTH_LIBPQ_H
#define __DBAUTH_LIBPQ_H
#ifdef __cplusplus
extern "C" {
#endif
    // Connect to the authenticator service to retrieve the DB user.
    // @param authenticatoraddress the address:port of the authenticator service.
    // @param identity identity string of the service.
    // @param dbuser pointer to memory where retreived dbuser will be saved.
    // @return 1 if the call was successful and 0 otherwise.
int dbauth_get_user(char *authenticatoraddress, char *identity, char **dbuser);
    // Connect to the authenticator service to retrieve the DB challenge MD5 hash.
    // @param authenticatoraddress the address:port of the authenticator service.
    // @param identity identity string of the service.
    // @param md5Salt point to the start of the md5Salt.
    // @param saltlen length of md5 salt in bytes.
    // @param hashbuffer buffer where result hash will be stored. has to be at least 36 bytes long
    // @return 1 if the call was successful and 0 otherwise.
int dbauth_get_hash(char *authenticatoraddress, char *identity, char *md5salt, int saltlen, char *hashbuffer);
#ifdef __cplusplus
}
#endif
#endif
