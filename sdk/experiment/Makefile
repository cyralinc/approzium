.PHONY: clean


CFLAGS   = 
INCLUDES =
LIBS     = 
LD_PATH  =

libhook.so: hook.c
		gcc $(CFLAGS) $(INCLUDES) -L$(LD_PATH) -Wl,--export-dynamic -fPIC -ggdb -shared -o $@ $^ -ldl $(LIBS)

clean:
		rm -f libhook.so
