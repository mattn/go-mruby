package main

import "github.com/mattn/go-mruby"

func main() {
	mrb := mruby.New()
	for _, i := range mrb.Eval(`[1,2,3].map {|x| x + 1}`).([]interface{}) {
		println(i.(int32)) // 2 3 4
	}
}
