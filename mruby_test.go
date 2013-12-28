package mruby

import (
	"reflect"
	"testing"
)

func TestMRuby(t *testing.T) {
	mruby := New()
	defer mruby.Close()
	s := "1 + 1"
	o := mruby.Eval(s)
	v, ok := o.(int32)
	if !ok {
		t.Errorf("Expected `%s` to yield an int32, got %v", s, reflect.TypeOf(o))
	} else if v != 2 {
		t.Errorf("Expected `%s` to equal 2, got %#v", s, v)
	}
}
