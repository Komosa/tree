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
	return ins(&t.root, k)
}

func ins(root **node, k Key) bool {
	if *root == nil {
		*root = &node{key: k}
		return true
	}

	x := *root
	eq, i := cmp(x.key, k)
	if eq {
		return false
	}

	return ins(&x.c[i], k)
}
