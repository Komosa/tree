package scapegoat

import (
	"math/rand"
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

func TestDelete(t *testing.T) {
	const (
		zero = Key(iota)
		one
		two
		three
	)

	tree := New(.65)
	eq(t, false, tree.Del(one))

	tree.Ins(one)
	eq(t, false, tree.Del(two))
	eq(t, true, tree.Del(one))
	eq(t, false, tree.Del(one))
	eq(t, false, tree.Del(two))

	tree.Ins(one)
	tree.Ins(two)
	eq(t, false, tree.Del(zero))
	eq(t, false, tree.Del(three))
	eq(t, true, tree.Del(one))
	eq(t, false, tree.Del(one))
	eq(t, true, tree.Del(two))

	tree.Ins(one)
	tree.Ins(two)
	tree.Ins(zero)
	eq(t, false, tree.Del(three))
	eq(t, true, tree.Del(one))
	eq(t, false, tree.Del(one))
	eq(t, true, tree.Del(zero))
	eq(t, false, tree.Del(zero))
	eq(t, true, tree.Del(two))
	eq(t, false, tree.Del(two))

	tree.Ins(one)
	tree.Ins(two)
	tree.Ins(zero)
	eq(t, true, tree.Del(one))
	eq(t, true, tree.Del(two))
	eq(t, true, tree.Del(zero))

	tree.Ins(one)
	tree.Ins(two)
	tree.Ins(zero)
	tree.Ins(three)
	eq(t, true, tree.Del(one))
	eq(t, true, tree.Del(three))
	eq(t, true, tree.Del(two))
	eq(t, true, tree.Del(zero))

	tree.Ins(one)
	tree.Ins(zero)
	tree.Ins(three)
	tree.Ins(two)
	eq(t, true, tree.Del(one))
	eq(t, true, tree.Del(three))
	eq(t, true, tree.Del(two))
	eq(t, true, tree.Del(zero))
}

func heightBalanced(x *node) bool {
	min, max := -1, -1
	var f func(*node, int)
	f = func(x *node, curr int) {
		if x == nil {
			if min == -1 || curr < min {
				min = curr
			}
			if curr > max {
				max = curr
			}
			return
		}
		f(x.c[0], curr+1)
		f(x.c[1], curr+1)
	}
	f(x, 0)
	return min+1 >= max
}

func TestRebalance(t *testing.T) {
	var expected []int

	for n := 1; n < 100; n++ {
		expected = append(expected, n-1)

		for tc := 1; tc <= n; tc++ {
			p := rand.Perm(n)
			tree := New(.65)

			for _, x := range p {
				tree.Ins(Key(x))
			}

			tree.root = rebalance(tree.root, n)

			var seq []int
			for it := tree.First(); it.Ok(); it = it.Next() {
				seq = append(seq, int(it.Key()))
			}

			eq(t, seq, expected, p)
			eq(t, heightBalanced(tree.root), true)
		}
	}
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
