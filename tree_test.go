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
		key byte
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

	var seq []byte
	for it := tree.First(); it.Ok(); it.Inc() {
		seq = append(seq, it.Key())
	}

	keys := []byte{1, 2, 3, 4, 6, 7}
	eq(t, seq, keys)

	for _, k := range keys {
		eq(t, tree.Exist(k), true, k)
		eq(t, tree.Exist(-k), false, -k)
	}
	eq(t, tree.Exist(0), false)
	eq(t, tree.Exist(5), false)
	eq(t, tree.Exist(8), false)
	eq(t, tree.Exist(9), false)
}

func TestDelete(t *testing.T) {
	const (
		zero = iota
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
	var expected []byte

	for n := 1; n < 100; n++ {
		expected = append(expected, byte(n-1))

		for tc := 1; tc <= n; tc++ {
			p := rand.Perm(n)
			tree := New(rand.Float64()/2 + 0.5)

			for _, x := range p {
				tree.Ins(byte(x))
			}

			tree.root = rebalance(tree.root, n)

			var seq []byte
			for it := tree.First(); it.Ok(); it.Inc() {
				seq = append(seq, byte(it.Key()))
			}

			eq(t, seq, expected, p)
			eq(t, heightBalanced(tree.root), true)
		}
	}
}

func TestIter(t *testing.T) {
	p := rand.Perm(100)
	tree := New(.6)
	for _, x := range p {
		tree.Ins(byte(x))
	}
	for it, exp := tree.Last(), byte(99); it.Ok(); it.Dec() {
		eq(t, it.Key(), exp)
		exp--
	}
	it := tree.First()
	it.Inc()
	it.Dec()
	eq(t, it.Key(), byte(0))

	for it.Key() != 41 {
		it.Inc()
	}
	for i := 0; i < 10; i++ {
		eq(t, it.Key(), byte(41))
		it.Inc()
		eq(t, it.Key(), byte(42))
		it.Dec()
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
