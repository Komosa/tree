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

func rebalance(x *node, subsize int) *node {
	fullCnt := 1
	for fullCnt*2+1 < subsize {
		fullCnt += fullCnt + 1
	}

	evenLeft := subsize - fullCnt
	even := make([]*node, 0, evenLeft+(evenLeft&1))
	odd := make([]*node, 0, fullCnt)

	for it := first(x); it.Ok(); {
		x := it[len(it)-1]

		if evenLeft > 0 && len(even) == len(odd) {
			even = append(even, x)
			evenLeft--
		} else {
			odd = append(odd, x)
		}

		it = goleft(x.c[1], it[:len(it)-1])
	}

	for _, x := range even {
		x.c[0] = nil
		x.c[1] = nil
	}
	if len(even)&1 == 1 {
		even = append(even, nil)
	}

	for i, x := range odd {
		j := ((i ^ (i + 1)) & (i + 1)) >> 1
		if j == 0 {
			if i < len(even) {
				x.c[0] = even[i]
				x.c[1] = even[i+1]
			} else {
				x.c[0] = nil
				x.c[1] = nil
			}
		} else {
			x.c[0] = odd[i-j]
			x.c[1] = odd[i+j]
		}
	}

	return odd[fullCnt/2]
}

type iterator []*node

func (t Tree) First() iterator {
	return first(t.root)
}

func (it iterator) Ok() bool {
	return len(it) > 0
}

// User must be sure that Ok() is true before call.
func (it iterator) Next() iterator {
	x := it[len(it)-1].c[1]
	it = it[:len(it)-1]
	return goleft(x, it)
}

// User must be sure that Ok() is true before call.
func (it iterator) Key() Key {
	return it[len(it)-1].key
}

// Make iterator over subtree rooted at given node.
func first(x *node) iterator {
	return goleft(x, iterator{})
}

func goleft(x *node, it iterator) iterator {
	for ; x != nil; x = x.c[0] {
		it = append(it, x)
	}
	return it
}
