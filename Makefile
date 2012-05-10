all : libmruby.so
	go build -x .

libmruby.so :
	(cd mruby && make)
	gcc -Wl,--export-all-symbols -shared -o libmruby.so mruby/lib/libmruby.a
