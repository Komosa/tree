package scapegoat

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestInsert(t *testing.T) {
	tree := New(.65)

	tcs := []struct {
		key Key
		ok  bool
	}{
		{1, true},
		{1, false},
		{2, true},
		{3, true},
		{4, true},
		{2, false},
		{1, false},
		{2, false},
		{3, false},
		{7, true},
		{7, false},
		{6, true},
	}
	for _, tc := range tcs {
		inserted := tree.Ins(tc.key)
		eq(t, inserted, tc.ok, tc)
	}

	var seq []Key
	for it := tree.First(); it.Ok(); it = it.Next() {
		seq = append(seq, it.Key())
	}

	keys := []Key{1, 2, 3, 4, 6, 7}
	eq(t, seq, keys)

	for _, k := range keys {
		eq(t, tree.Exist(k), true, k)
		eq(t, tree.Exist(-k), false, -k)
	}
	eq(t, tree.Exist(Key(0)), false)
	eq(t, tree.Exist(Key(5)), false)
	eq(t, tree.Exist(Key(8)), false)
	eq(t, tree.Exist(Key(9)), false)
}

func eq(tb testing.TB, act, exp interface{}, info ...interface{}) {
	if !reflect.DeepEqual(act, exp) {
		_, file, line, _ := runtime.Caller(1)
		if len(info) > 0 {
			tb.Errorf("%s:%d: got:%#v but exp:%#v with:%v\n", filepath.Base(file), line, act, exp, info)
		} else {
			tb.Errorf("%s:%d: got:%#v but exp:%#v\n", filepath.Base(file), line, act, exp)
		}
	}
}
