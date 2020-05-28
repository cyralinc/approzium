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

void print_source(char *p, int offset, int nbytes);
static int unprotect_page(void* addr);

char *hook2() {
    fprintf(stderr, "hook2\n");
    char *p = (char *) &PQfireResultCreateEvents;
    char target[] = {0xba, 0x04, 0x00, 0x00, 0x00, 0x4c};
    int max_offset = 10000000;
    char *target_p = NULL;
    for (int i=0; i < max_offset; i--) {
        int j = 0;
        for (; j < sizeof(target); j++) {
            if (p[i+j] != target[j])
                break;
        }
        if (j == sizeof(target)) {
            target_p = (char *) p+i;
            fprintf(stderr, "target_p=%p\n", target_p);
            print_source(target_p, -16, 32);
        }
    }
    if (target_p == NULL) {
        fprintf(stderr, "failed to find target\n");
        return NULL;
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
    void *p1 = &PQconnectStart;
    void *p2 = &PQlibVersion;
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
