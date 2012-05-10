package mruby

/*
#include "mruby/include/mruby.h"
#include "mruby/include/mruby/proc.h"
#include "mruby/include/mruby/data.h"
#undef INCLUDE_ENCODING
#include "mruby/include/mruby/string.h"
#include "mruby/include/mruby/khash.h"
#include "mruby/include/mruby/hash.h"
#include "mruby/include/mruby/array.h"
#include "mruby/include/mruby/class.h"
#include "mruby/include/mruby/object.h"
#include "mruby/include/mruby/variable.h"
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

#cgo LDFLAGS: ./libmruby.dll.a
*/
import "C"
import "unsafe"

type MRuby struct {
	mrb *C.mrb_state
}

func New() *MRuby {
	return &MRuby { C.mrb_open() }
}

func mruby2go(mrb *C.mrb_state, o C.mrb_value) interface{} {
	println(o.tt)
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

func (m *MRuby) Run(code string) {
	c := C.CString(code)
	defer C.free(unsafe.Pointer(c))
	p := C.mrb_parse_string(m.mrb, c)
	n := C.mrb_generate_code(m.mrb, p.tree);
	C.mrb_pool_close((*C.mrb_pool)(p.pool))
	if n >= 0 {
		C.mrb_run(
			m.mrb,
			C.mrb_proc_new(m.mrb, (*C.mrb_irep)(C._get_irep(m.mrb, n))),
			C.mrb_top_self(m.mrb))
	}

	if m.mrb.exc != nil {
        C.mrb_p(m.mrb, C.mrb_obj_value(unsafe.Pointer(m.mrb.exc)))
	}
}

func (m *MRuby) Eval(code string) interface{} {
	c := C.CString(code)
	defer C.free(unsafe.Pointer(c))
	p := C.mrb_parse_string(m.mrb, c)
	n := C.mrb_generate_code(m.mrb, p.tree);
	C.mrb_pool_close((*C.mrb_pool)(p.pool))
	if n >= 0 {
		C.mrb_run(
			m.mrb,
			C.mrb_proc_new(m.mrb, (*C.mrb_irep)(C._get_irep(m.mrb, n))),
			C.mrb_top_self(m.mrb))
	}

	if m.mrb.exc != nil {
		return mruby2go(m.mrb, C.mrb_obj_value(unsafe.Pointer(m.mrb.exc)))
	}
	return nil
}
