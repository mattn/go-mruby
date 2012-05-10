package main

import "github.com/mattn/go-mruby"

func main() {
	mrb := mruby.New()
	mrb.Eval(`
      [1,2,3].map {|x|
        puts x
      }
    `)
}
