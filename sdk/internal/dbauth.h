#include <ctype.h>
#include <stdarg.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <stdlib.h>
#include "authenticator/messages/authenticator.pb-c.h"
#include <protobuf-c-rpc/protobuf-c-rpc.h>
#include <stdio.h>
#include "postgres_fe.h"
#include "libpq-int.h"

/* PROTOTYPES */
char *dbauth_get_hashed_password(PGconn *conn, char *md5Salt);
void handle_query_response(const  Dbauth__Authenticator__Messages__AuthenticateResponse *result, void *hashedpassword);

void handle_query_response(const  Dbauth__Authenticator__Messages__AuthenticateResponse *result, void *hashedpassword) {
	if (result == NULL) {
		fprintf(stderr, "Error processing request\n");
		*(char **) hashedpassword = NULL;
		return;
	}

	if (result->status != DBAUTH__AUTHENTICATOR__MESSAGES__AUTHENTICATE_RESPONSE__STATUS__SUCCESS) {
		fprintf(stderr, "Authenticate request failed\n");
		*(char **) hashedpassword = NULL;
		return;
	}
	*(char **) hashedpassword = result->credentials->hashedpassword;
}

char *dbauth_get_hashed_password(PGconn *conn, char *md5Salt) {
	ProtobufCService *service;
	ProtobufC_RPC_Client *client;
	char *unset; 
	char *hashedpassword;

	fprintf(stderr, "dbauth_get_hashed_password\n");
	fprintf(stderr, "MD5 Salt: %u\n", *(uint32_t *)md5Salt);
	service = protobuf_c_rpc_client_new(PROTOBUF_C_RPC_ADDRESS_TCP, "172.18.0.2:1234", &dbauth__authenticator__messages__authenticator__descriptor, NULL);
	if (service == NULL) {
		fprintf(stderr, "could not create service\n");
		return NULL;
	}
	client = (ProtobufC_RPC_Client *) service;
	fprintf(stderr, "connecting...\n");
	while (!protobuf_c_rpc_client_is_connected(client))
		protobuf_c_rpc_dispatch_run(protobuf_c_rpc_dispatch_default());
	fprintf(stderr, "done.\n");

	Dbauth__Authenticator__Messages__AuthenticateRequest query = DBAUTH__AUTHENTICATOR__MESSAGES__AUTHENTICATE_REQUEST__INIT;
	query.identity = "diobrando";
	query.salt = md5Salt;
	unset = "unset";
	hashedpassword = unset;
	dbauth__authenticator__messages__authenticator__authenticate(service, &query, handle_query_response, &hashedpassword); 
	while (hashedpassword == unset)
		protobuf_c_rpc_dispatch_run(protobuf_c_rpc_dispatch_default());
	return hashedpassword;
}
