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
extern void PQencryptPassword();

void print_source(char *p, int offset, int nbytes);
static int unprotect_page(void* addr);


python_md5_func *callback = NULL;

void register_python_callback(python_md5_func *_callback) {
    callback = _callback;
}

bool our_md5_encrypt(const char *passwd, const char *salt, size_t salt_len,
        char *buf) {
    const unsigned char *usalt = (const unsigned char *)salt;
    fprintf(stderr, "libhook(our_md5_encrypt): salt (%ld bytes) is %02x,%02x,%02x,%02x\n",
            salt_len, usalt[0], usalt[1], usalt[2], usalt[3]);
    if (callback != NULL) {
        char *salt_null_terminated = (char *) malloc(salt_len+1);
        memcpy(salt_null_terminated, salt, salt_len);
        salt_null_terminated[salt_len] = '\0';
        char *received_hash = (*callback)(salt_null_terminated);
        free(salt_null_terminated);
    fprintf(stderr, "libhook(our_md5_encrypt): received hash %s\n", received_hash);
    }
    buf[0] = 'd';
    buf[1] = 'b';
    buf[2] = 'a';
    for (char i =3; i < 35; i++) {
        buf[i] = i;
    }
    buf[35] = '\0';
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

char *hook1() {
    // this hook shows that we can overrite the machine instructions of a
    // function to change its return value. in this example, the (exported) 
    // function PQProtocolVersion is changed so that its returns 420, which is
    // not a possible protocol number. this works. the next challenge is to change
    // that is not public.

    char *p1 = (char *) &PQprotocolVersion;
    int index = 0x14;
    char overwrite[] = {0xb8, 0xa4, 0x01, 0x00, 0x00, 0x90, 0x90};
	if(unprotect_page(p1)) {
		fprintf(stderr, "Could not unprotect mem: %p\n", p1);
		return NULL;
	}
    for (int i=0; i<sizeof(overwrite); i++) {
        p1[index+i] = overwrite[i];
    }
    return NULL;
}

void print_source(char *p, int offset, int nbytes) {
    for (int i=offset; i < offset+nbytes; i++) {
        if (i == 0)
            fprintf(stderr, "|");
        fprintf(stderr, "%02x ", (unsigned char) p[i]);
    }
    fprintf(stderr, "\n");
}


int sanity_check1() {
    // this checks that at least for these two functions, the relative addresses
    // do not change at different runtimes. if they do, that means that the
    // dynamic linking is doing more than a simple memory map, or that there is
    // something else funky going on.
    // this used to be true for apt libpq-dev but isnot for compiled libpq
    void *p1 = &PQconnectStart;
    void *p2 = &PQlibVersion;
    fprintf(stderr, "%p\n ", p2);
    return (p2-p1) == 18640;
}

/* TODO: make less ugly! 
 *       there's got to be a nicer way to do this! */
#pragma pack(push, 1)
static struct { 
	char push_rax; 
	char mov_rax[2];
	char addr[8];
	char jmp_rax[2]; } 
jump_asm = {
	.push_rax = 0x50,
	.mov_rax  = {0x48, 0xb8},
	.jmp_rax  = {0xff, 0xe0} };
#pragma pack(pop)

static int unprotect_page(void* addr) {
	int pagesize = sysconf(_SC_PAGE_SIZE);
	int pagemask = ~(pagesize -1);
	char* page = (char *)((size_t)addr & pagemask);
	return mprotect(page, pagesize, PROT_READ | PROT_WRITE | PROT_EXEC);
}

int hook_function(void* target, void* replace) {
	int count;

	if(unprotect_page(replace)) {
		fprintf(stderr, "Could not unprotect replace mem: %p\n", replace);
		return 1;
	}

	if(unprotect_page(target)) {
		fprintf(stderr, "Could not unprotect target mem: %p\n", target);
		return 1;
	}

	/* find the NOP */
	for(count = 0; count < 255 && ((unsigned char*)replace)[count] != 0x90; ++count);

	if(count == 255) {
		fprintf(stderr, "Couldn't find the NOP.\n");
		return 1;
	}

	/* shift everything down one */
	memmove(replace+1, replace, count);

	/* add in `pop %rax` */
	*((unsigned char *)replace) = 0x58;

	/* set up the address */
	memcpy(jump_asm.addr, &replace, sizeof (void *));

	/* smash the target function */
	memcpy(target, &jump_asm, sizeof jump_asm);

	return 0;
}
