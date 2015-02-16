package scapegoat

type Key int

func Less(a, b Key) bool {
	return a < b
}
func Eq(a, b Key) bool {
	return a == b
}
func cmp(x *node, k Key) (eq bool, idx int) {
	if Less(x.key, k) {
		return false, 1
	}
	return Eq(x.key, k), 0
}

type node struct {
	key Key
	c   [2]*node
}

type Tree struct {
	root *node
	alfa float64
}

func New(alfa float64) *Tree {
	return &Tree{alfa: alfa}
}

func (t *Tree) Ins(k Key) bool {
	for x, p := t.root, &t.root; ; {
		if x == nil {
			*p = &node{key: k}
			return true
		}

		eq, i := cmp(x, k)
		if eq {
			return false
		}

		p = &x.c[i]
		x = x.c[i]
	}
}

func (t Tree) Exist(k Key) bool {
	for x := t.root; x != nil; {
		eq, i := cmp(x, k)
		if eq {
			return true
		}
		x = x.c[i]
	}
	return false
}

func (t *Tree) Del(k Key) bool {
	p := &t.root
	x := t.root
	for x != nil {
		eq, i := cmp(x, k)
		if eq {
			break
		}
		p = &x.c[i]
		x = x.c[i]
	}

	if x == nil {
		return false
	}

	if x.c[0] != nil && x.c[1] != nil {
		y := x.c[1]
		p := &x.c[1]
		for y.c[0] != nil {
			p = &y.c[0]
			y = y.c[0]
		}

		*p = y.c[1] // it is fine both for nil and non-nil

		x.key = y.key
		return true
	}

	if x.c[0] == nil {
		*p = x.c[1] // it is fine also for nil
	} else {
		*p = x.c[0]
	}
	return true
}

type iterator []*node

func (t Tree) First() iterator {
	var s iterator

	for x := t.root; x != nil; x = x.c[0] {
		s = append(s, x)
	}
	return s
}

func (it iterator) Ok() bool {
	return len(it) > 0
}

// User must be sure that Ok() is true before call.
func (it iterator) Next() iterator {
	x := it[len(it)-1].c[1]
	it = it[:len(it)-1]
	for ; x != nil; x = x.c[0] {
		it = append(it, x)
	}
	return it
}

// User must be sure that Ok() is true before call.
func (it iterator) Key() Key {
	return it[len(it)-1].key
}
