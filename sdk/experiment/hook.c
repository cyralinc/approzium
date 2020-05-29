#include <stdlib.h>
#include <stdbool.h>
#include <stdio.h>
#include <sys/mman.h>
#include <unistd.h>
#include <stdint.h>
#include <string.h>

#include "hook.h"

extern void PQconnectStart();
extern void PQlibVersion();
extern void PQprotocolVersion();
extern void PQfireResultCreateEvents();
extern char *PQencryptPassword(const char *passwd, const char *user);

void print_source(char *p, int offset, int nbytes);
static int unprotect_page(void* addr);


python_md5_func *callback = NULL;

void register_python_callback(python_md5_func *_callback) {
    callback = _callback;
}

bool our_md5_encrypt(const char *passwd, const char *salt, size_t salt_len,
        char *buf) {
    char *salt_string;
    char *hash;
    salt_string = (char *) malloc(salt_len+1);
    memcpy(salt_string, salt, salt_len);
    salt_string[salt_len] = '\0';
    fprintf(stderr, "libhook(our_md5_encrypt): salt (%ld bytes) is ", salt_len);
    for (int i=0; i < salt_len; i++) {
        fprintf(stderr, "%02x", (unsigned char) salt_string[i]);
    }
    fprintf(stderr, "\n");
    if (callback != NULL) {
        hash = (*callback)(salt_string);
        if (strlen(hash) > 0) {
            fprintf(stderr, "libhook(our_md5_encrypt): received hash %s\n", hash);
        } else {
            fprintf(stderr, "libhook(our_md5_encrypt): received empty hash, reverting to libpq hash\n");
            hash = PQencryptPassword(passwd, salt_string);
            if (!hash)
                return false;
        }
    }
    memcpy(buf, hash, strlen(hash));
    free(hash);
    return true;
}

char *hook2() {
    // this hook overwrides the 2nd call to pg_md5_encrypt in the pg_password_sendauth
    // function. whatever this send call returns is considered the final hash that is sent
    // as the challenge response. therefore, we can take over that function and return
    // a hash that we computed/received from somewhere else. this implementation works by
    // calculating a new offset value for the `call` assembly instruction. this offset value
    // points to our own `our_md5_encrypt` function instead of `pg_md5_encrypt`
    fprintf(stderr, "libhook(hook2)\n");
    char *start = (char *) &PQfireResultCreateEvents;
    char *end = (char *) &PQencryptPassword;
    if (start > end) {
        char *temp = end;
        end = start;
        start = temp;
    }
    char target[] = {0xe8, 0x1b, 0x3a, 0x00, 0x00};
    char *target_p = NULL;
    for (char *i=start; i < end; i++) {
        int j = 0;
        for (; j < sizeof(target); j++) {
            if (i[j] != target[j])
                break;
        }
        if (j == sizeof(target)) {
            target_p = (char *) i;
        }
    }
    if (target_p == NULL) {
        fprintf(stderr, "failed to find target\n");
        return NULL;
    }
    // patch address
    int64_t new_offset = (int64_t) &our_md5_encrypt - ((int64_t) target_p+5);
    char *call_offset = target_p + 1;
	if(unprotect_page(call_offset)) {
		fprintf(stderr, "Could not unprotect mem: %p\n", call_offset);
		return NULL;
	}
    for (char i=0; i < 4; i++) {
        unsigned char b = (unsigned char) (new_offset >> (8*i));
        call_offset[i] = b;
    }
    return NULL;
}

static int unprotect_page(void* addr) {
    int pagesize = sysconf(_SC_PAGE_SIZE);
    int pagemask = ~(pagesize -1);
    char* page = (char *)((size_t)addr & pagemask);
    return mprotect(page, pagesize, PROT_READ | PROT_WRITE | PROT_EXEC);
}
