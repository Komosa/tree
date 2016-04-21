package tree

import (
	"errors"
	"math"
)

var ErrMovePastEnd = errors.New("tree: iterator moved past to end")

// key → byte

// cmp should compare two keys
// first return value should indicate if both keys are equal
// second should note into which child we should descend in order to look
// for key k form this node
func (n node) cmp(k byte) (bool, int) {
	idx := 0
	if n.key < k {
		idx = 1
	}
	return k == n.key, idx
}

type node struct {
	key byte
	c   [2]*node
}

type Tree struct {
	root    *node
	alfa    float64
	size    int
	maxsize int
}

type Iter struct {
	stack []*node
}

func New(alfa float64) *Tree {
	return &Tree{alfa: alfa}
}

func (t *Tree) Ins(k byte) bool {
	t.size++
	if t.root == nil {
		t.root = &node{key: k}
		return true
	}

	var path []*node
	for x := t.root; x != nil; {
		eq, i := x.cmp(k)
		if eq {
			t.size--
			return false
		}

		path = append(path, x)
		x = x.c[i]
	}

	x := path[len(path)-1]
	_, j := x.cmp(k)
	x.c[j] = &node{key: k}
	if t.size > t.maxsize {
		t.maxsize = t.size
	}

	if pow(1/t.alfa, len(path)) <= float64(t.size) {
		return true
	}

	var scapegoat int
	ssize := t.size
	for i, childsize := len(path)-1, 1; i > 0; i-- {
		x = path[i]
		ochildsize := subsize(x.c[j^1])
		currsize := childsize + ochildsize + 1
		if math.Max(float64(childsize), float64(ochildsize)) > float64(currsize)*t.alfa {
			scapegoat = i
			ssize = currsize
			break
		}
		childsize = currsize
		_, j = path[i-1].cmp(k)
	}

	x = rebalance(path[scapegoat], ssize)
	if scapegoat == 0 {
		t.root = x
		t.maxsize = ssize
	} else {
		_, i := path[scapegoat-1].cmp(k)
		path[scapegoat-1].c[i] = x
	}

	return true
}

func (t Tree) Exist(k byte) bool {
	for x := t.root; x != nil; {
		eq, i := x.cmp(k)
		if eq {
			return true
		}
		x = x.c[i]
	}
	return false
}

func (t *Tree) Del(k byte) bool {
	p := &t.root
	x := t.root
	for x != nil {
		eq, i := x.cmp(k)
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
		p = &x.c[1]
		for y.c[0] != nil {
			p = &y.c[0]
			y = y.c[0]
		}

		*p = y.c[1] // it is fine both for nil and non-nil

		x.key = y.key
	} else if x.c[0] == nil {
		*p = x.c[1] // it is fine also for nil
	} else {
		*p = x.c[0]
	}

	t.size--
	if t.size > 0 && float64(t.size) <= t.alfa*float64(t.maxsize) {
		t.maxsize = t.size
		t.root = rebalance(t.root, t.size)
	}
	return true
}

func rebalance(x *node, subsize int) *node {
	// we will construct full tree plus some (even) nodes appended to its leafs
	fullCnt := 1
	for fullCnt <= subsize {
		fullCnt += fullCnt + 1
	}
	fullCnt /= 2
	evenLeft := subsize - fullCnt

	d := make([]*node, subsize)
	even := d[:0]
	odd := d[evenLeft:evenLeft]

	for it := first(0, &Iter{stack: []*node{x}}); it.Ok(); it.Inc() {
		x := it.top()

		if evenLeft > 0 && len(even) == len(odd) {
			even = append(even, x)
			evenLeft--
		} else {
			odd = append(odd, x)
		}
	}

	for _, x := range d {
		x.c[0] = nil
		x.c[1] = nil
	}
	for i, x := range even {
		odd[i&^1].c[i&1] = x
	}
	for i := 1; i < len(odd); i += 2 {
		j := ((i ^ (i + 1)) & (i + 1)) >> 1
		odd[i].c[0] = odd[i-j]
		odd[i].c[1] = odd[i+j]
	}

	return odd[fullCnt/2]
}

func (it *Iter) top() *node { return it.stack[len(it.stack)-1] }
func (it *Iter) pop() *node { x := it.top(); it.stack = it.stack[:len(it.stack)-1]; return x }

func first(dir int, it *Iter) *Iter {
	x := it.pop()
	if x == nil {
		return it
	}

	for x.c[dir] != nil {
		it.stack = append(it.stack, x)
		x = x.c[dir]
	}
	it.stack = append(it.stack, x)
	return it
}

func move(dir int, it *Iter) error {
	x := it.pop()
	if x == nil {
		return ErrMovePastEnd
	}
	if x.c[dir^1] == nil {
		for {
			if len(it.stack) == 0 {
				return nil
			}

			if it.top().c[dir] == x {
				x = it.pop()
				it.stack = append(it.stack, x)
				return nil
			}

			x = it.pop()
		}
	}
	it.stack = append(it.stack, x)
	x = x.c[dir^1]
	it.stack = append(it.stack, x)
	first(dir, it)
	return nil
}

func (t Tree) First() *Iter { return first(0, &Iter{stack: []*node{t.root}}) }
func (t Tree) Last() *Iter  { return first(1, &Iter{stack: []*node{t.root}}) }
func (it *Iter) Inc() error { return move(0, it) }
func (it *Iter) Dec() error { return move(1, it) }
func (it *Iter) Ok() bool   { return len(it.stack) > 0 }

// User must be sure that Ok() is true before call.
func (it *Iter) Key() byte { return it.top().key }

func subsize(x *node) int {
	if x == nil {
		return 0
	}
	return 1 + subsize(x.c[0]) + subsize(x.c[1])
}

func pow(x float64, n int) float64 {
	res := 1.0
	for n > 0 {
		if n&1 == 1 {
			res *= x
		}
		x *= x
		n >>= 1
	}
	return res
}
