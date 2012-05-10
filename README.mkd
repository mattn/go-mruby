# go-mruby

go-mruby make interface to embed mruby into go.

## Install

```
git submodule init
git submodule update
make
cd example
go build -x .
LD_LIBRARY_PATH=.. ./example
```

On windows, use Makefile.w32

```
mingw32-make -f Makefile.w32
copy mruby.dll example
cd example
go build -x .
example.exe
```


## Usage

```go
package main

import "github.com/mattn/go-mruby"

func main() {
	mrb := mruby.New()
	defer mrb.Close()

	println(mrb.Eval(`"hello " + ARGV[0]`, "mruby").(string))

	for _, i := range mrb.Eval(`ARGV.map {|x| x + 1}`, 1, 2, 3).([]interface{}) {
		println(i.(int32)) // 2 3 4
	}
}
```

## License

MIT

## Author

* Yasuhiro Matsumoto
