all : libmruby.so
	go build -x .

libmruby.so : mruby/lib/libmruby.a
	ld --whole-archive -shared -o libmruby.so mruby/lib/libmruby.a

mruby/lib/libmruby.a :
	(cd mruby && make)

clean :
	(cd mruby && make clean)
	go clean .
