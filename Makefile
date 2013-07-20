all : mruby/lib/libmruby.a
	go build -x .

mruby/lib/libmruby.a :
	(cd mruby && make)

clean :
	(cd mruby && make clean)
	go clean .
