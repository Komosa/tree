// +build graphviz

// use  go install -tags=graphviz
// or   go build   -tags=graphviz
// to support building this method

package scapegoat

import (
	"fmt"
	"io"
)

func Dot(t *Tree, w io.Writer) {
	io.WriteString(w, "digraph x {\nnode [shape=Mrecord]\n")
	var f func(x *node)
	limit := 1000
	f = func(x *node) {
		limit--
		if limit < 0 {
			return
		}
		fmt.Fprintf(w, "node%p [label=\"{%d|{<c0> L|<c1> r}}\"]\n", x, x.key)
		for i, y := range x.c {
			if y != nil {
				fmt.Fprintf(w, "node%p:<c%d> -> node%p\n", x, i, y)
				f(y)
			}
		}
	}
	if t.root != nil {
		f(t.root)
	}
	io.WriteString(w, "}\n")
}
