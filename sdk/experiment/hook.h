#ifndef HOOK_H
#define HOOK_H

#if !(__x86_64__)
	#error "This only works on x86_64"
#endif

typedef char *python_md5_func(char *);
void register_python_callback(python_md5_func *_callback);
char *hook2();
char *hook1();
int func1();
int hook_function(void* target, void* replace);

#endif
