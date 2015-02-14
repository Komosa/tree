package scapegoat

type Key int

func Less(a, b Key) bool {
	return a < b
}
func Eq(a, b Key) bool {
	return a == b
}
func cmp(a, b Key) (eq bool, idx int) {
	if Less(a, b) {
		return false, 1
	}
	return Eq(a, b), 0
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

		eq, i := cmp(x.key, k)
		if eq {
			return false
		}

		p = &x.c[i]
		x = x.c[i]
	}
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
