package mruby

/*
#include <mruby.h>
#include <mruby/proc.h>
#include <mruby/data.h>
#undef INCLUDE_ENCODING
#include <mruby/string.h>
#include <mruby/khash.h>
#include <mruby/hash.h>
#include <mruby/array.h>
#include <mruby/class.h>
#include <mruby/object.h>
#include <mruby/variable.h>
#include "mruby/src/compile.h"

extern khint_t mrb_hash_ht_hash_func(mrb_state *mrb, mrb_value key);
extern khint_t mrb_hash_ht_hash_equal(mrb_state *mrb, mrb_value a, mrb_value b);

KHASH_INIT(ht, mrb_value, mrb_value, 1, mrb_hash_ht_hash_func, mrb_hash_ht_hash_equal);

static mrb_value
mrb_hash_keys(mrb_state *mrb, mrb_value hash) {
  khash_t(ht) *h = RHASH_TBL(hash);
  khiter_t k;
  mrb_value ary = mrb_ary_new(mrb);

  if (!h) return ary;
  for (k = kh_begin(h); k != kh_end(h); k++) {
    if (kh_exist(h, k)) {
      mrb_value v = kh_key(h,k);
      if ( !mrb_special_const_p(v) )
        v = mrb_obj_dup(mrb, v);
      mrb_ary_push(mrb, ary, v);
    }
  }
  return ary;
}

static int _RARRAY_LEN(mrb_value a) { return (RARRAY(a)->len); }
static int _mrb_fixnum(mrb_value o) { return (int) mrb_fixnum(o); }
static float _mrb_float(mrb_value o) { return (float) mrb_float(o); }
static struct mrb_irep* _get_irep(mrb_state *mrb, int n) { return mrb->irep[n]; }
//static mrb_value _get_result(mrb_state *mrb, int n) { return mrb->stack[mrb->irep[n]->nlocals]; }

#cgo CFLAGS: -I. -Imruby/include
#cgo linux LDFLAGS: -L. libmruby.so -lm
#cgo windows LDFLAGS: ./libmruby.dll.a
*/
import "C"
import "unsafe"
import "reflect"

type MRuby struct {
	mrb *C.mrb_state
}

func New() *MRuby {
	return &MRuby { C.mrb_open() }
}

func go2mruby(mrb *C.mrb_state, o interface{}) C.mrb_value {
	v := reflect.ValueOf(o)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return C.mrb_fixnum_value(C.mrb_int(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return C.mrb_fixnum_value(C.mrb_int(v.Uint()))
	case reflect.Float32, reflect.Float64:
		return C.mrb_float_value((C.mrb_float)(v.Float()))
	case reflect.Complex64, reflect.Complex128:
		return C.mrb_float_value((C.mrb_float)(v.Float()))
	case reflect.String:
		ptr := C.CString(v.String())
		return C.mrb_str_new(mrb, ptr, C.strlen(ptr))
	case reflect.Bool:
		if v.Bool() {
			return C.mrb_true_value()
		}
		return C.mrb_false_value()
	case reflect.Array, reflect.Slice:
		ary := C.mrb_ary_new(mrb)
		for i := 0; i < v.Len(); i++ {
			C.mrb_ary_push(mrb, ary, go2mruby(mrb, v.Index(i)))
		}
		return ary
	case reflect.Map:
		hash := C.mrb_hash_new(mrb, 32)
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			C.mrb_hash_set(mrb, hash, go2mruby(mrb, key.String()), go2mruby(mrb, val))
		}
		return hash
	}
	return C.mrb_nil_value()
}

func mruby2go(mrb *C.mrb_state, o C.mrb_value) interface{} {
	switch o.tt {
	case C.MRB_TT_TRUE:
		return true
	case C.MRB_TT_FALSE:
		return false
	case C.MRB_TT_FLOAT:
		return float32(C._mrb_float(o))
	case C.MRB_TT_FIXNUM:
		return int32(C._mrb_fixnum(o))
	case C.MRB_TT_ARRAY:
		{
			var list []interface{}
			for i := 0; i < int(C._RARRAY_LEN(o)); i++ {
				list = append(list, mruby2go(mrb, C.mrb_ary_ref(mrb, o, C.mrb_int(i))))
			}
			return list
		}
	case C.MRB_TT_HASH:
		{
			hash := make(map[string]interface{})
			keys := C.mrb_hash_keys(mrb, o)
			for i := 0; i < int(C._RARRAY_LEN(keys)); i++ {
				key := C.mrb_ary_ref(mrb, keys, C.mrb_int(i))
				val := C.mrb_hash_get(mrb, o, key)
				hash[C.GoString(C.mrb_string_value_ptr(mrb, key))] = mruby2go(mrb, val)
			}
			return hash
		}
	case C.MRB_TT_STRING:
		return C.GoString(C.mrb_string_value_ptr(mrb, o))
	}
	return nil
}

func (m *MRuby) Run(code string, args ...interface{}) {
	c := C.CString(code)
	defer C.free(unsafe.Pointer(c))
	p := C.mrb_parse_string(m.mrb, c)
	n := C.mrb_generate_code(m.mrb, p.tree)
	C.mrb_pool_close((*C.mrb_pool)(p.pool))
	a := C.CString("ARGV")
	defer C.free(unsafe.Pointer(a))
	if n >= 0 {
		ARGV := C.mrb_ary_new(m.mrb)
		for i := 0; i < len(args); i++ {
			C.mrb_ary_push(m.mrb, ARGV, go2mruby(m.mrb, args[i]))
		}
		C.mrb_define_global_const(m.mrb, a, ARGV)
		C.mrb_run(
			m.mrb,
			C.mrb_proc_new(m.mrb, (*C.mrb_irep)(C._get_irep(m.mrb, n))),
			C.mrb_top_self(m.mrb))
	}

	if m.mrb.exc != nil {
        C.mrb_p(m.mrb, C.mrb_obj_value(unsafe.Pointer(m.mrb.exc)))
	}
}

func (m *MRuby) Eval(code string, args ...interface{}) interface{} {
	c := C.CString(code)
	defer C.free(unsafe.Pointer(c))
	p := C.mrb_parse_string(m.mrb, c)
	n := C.mrb_generate_code(m.mrb, p.tree)
	C.mrb_pool_close((*C.mrb_pool)(p.pool))
	a := C.CString("ARGV")
	defer C.free(unsafe.Pointer(a))
	if n >= 0 {
		ARGV := C.mrb_ary_new(m.mrb)
		for i := 0; i < len(args); i++ {
			C.mrb_ary_push(m.mrb, ARGV, go2mruby(m.mrb, args[i]))
		}
		C.mrb_define_global_const(m.mrb, a, ARGV)
		return mruby2go(m.mrb, C.mrb_run(
			m.mrb,
			C.mrb_proc_new(m.mrb, (*C.mrb_irep)(C._get_irep(m.mrb, n))),
			C.mrb_top_self(m.mrb)))
	}
	return C.mrb_nil_value()
}

func (m *MRuby) Close() {
	if m.mrb != nil {
		C.mrb_close(m.mrb)
	}
}
