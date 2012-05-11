all : libmruby.so
	go build -x .

libmruby.so : mruby/lib/libmruby.a
	gcc -Wl,-E -shared -o libmruby.so mruby/src/*.o mruby/mrblib/*.o

mruby/lib/libmruby.a :
	(cd mruby && make CFLAGS='-g -O3 -fPIC -shared')
