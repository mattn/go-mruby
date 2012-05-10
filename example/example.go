package main

import "github.com/mattn/go-mruby"

func main() {
	mrb := mruby.New()
	for i := range mrb.Eval(`[1,2,3].map {|x| x + 1}`).([]interface{}) {
		println(int(i))
	}
}
