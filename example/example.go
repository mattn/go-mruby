package main

import "github.com/mattn/go-mruby"

func main() {
	mrb := mruby.New()
	for _, i := range mrb.Eval(`ARGV.map {|x| x + 1}`, 1, 2, 3).([]interface{}) {
		println(i.(int32)) // 2 3 4
	}
	println(mrb.Eval(`"hello " + ARGV[0]`, "mruby").(string))
}
