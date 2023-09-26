package xmlparser

import (
	"bytes"
	"fmt"
	"strings"
)

type treeNode[T any] struct {
	data                     T
	index, first, last, next int
}

func (n treeNode[T]) Data() T {
	return n.data
}

func (n treeNode[T]) Index() int {
	return n.index
}

type tree[T any] struct {
	nodes []treeNode[T] //first node will be always last node
}

func (t *tree[T]) insertChild(parent *treeNode[T], n treeNode[T]) {
	n.index = len(t.nodes)
	t.nodes = append(t.nodes, n)
	if parent != nil {
		index := len(t.nodes) - 1
		if parent.first == -1 {
			//first node
			parent.first = index
		} else {
			//subsequent node
			t.nodes[parent.last].next = index
		}
		parent.last = index
	}
}

func (t *tree[T]) reset() {
	t.nodes = t.nodes[:0]
}

func (t *tree[T]) getMatchedChildrens(parent int, match func(T) bool, cb func(int, *treeNode[T])) {
	if parent < 0 || parent >= len(t.nodes) {
		return
	}
	for i := t.nodes[parent].first; i != -1; i = t.nodes[i].next {
		if match == nil || match(t.nodes[i].data) {
			cb(i, &t.nodes[i])
		}
	}
}

func (t *tree[T]) getChildrens(parent int, cb func(int, *treeNode[T])) {
	t.getMatchedChildrens(parent, nil, cb)
}

func (t *tree[T]) _get(parent int,
	path []string,
	match func(string, T) bool,
	cb func(int, *treeNode[T])) {

	t.getMatchedChildrens(parent,
		func(node T) bool {
			return match(path[0], node)
		},
		func(index int, node *treeNode[T]) {
			if len(path) == 1 {
				// last element
				cb(index, node)
			} else {
				t._get(index, path[1:], match, cb)
			}
		},
	)
}

func (t *tree[T]) get(path []string, match func(string, T) bool) (result *treeNode[T], found bool) {

	t._get(len(t.nodes)-1, path[:], match, func(_ int, node *treeNode[T]) {
		found = true
		result = node
	})
	return
}

/* Printing Function */
func (t *tree[T]) _print(buf *bytes.Buffer, index, indent int, f func(T) string) {
	buf.WriteByte('\n')
	buf.WriteString(strings.Repeat("\t", indent))
	buf.WriteByte('|')
	buf.WriteString(f(t.nodes[index].data))
	indent++
	for i := t.nodes[index].first; i != -1; i = t.nodes[i].next {
		t._print(buf, i, indent, f)
	}
}

func (t *tree[T]) print(f func(T) string) string {
	root := len(t.nodes) - 1
	buf := bytes.Buffer{}
	t._print(&buf, root, 0, f)
	return buf.String()
}

func (t *tree[T]) printRaw(f func(T) string) string {
	buf := bytes.Buffer{}
	for i, node := range t.nodes {
		buf.WriteString(fmt.Sprintf("\n%d:<%d,%d,%d>:%s", i, node.first, node.last, node.next, f(node.data)))
	}
	return buf.String()
}

/*
func (t *tree[T]) get(path []string, match func(string, T) bool) (*T, bool) {
	stackIndex := make([]int, 0, len(path)) //using pool
	pathIndex := 0
	i := len(t.nodes) - 1

	for i != -1 {
		//compare
		if match(path[pathIndex], t.nodes[i].data) {
			if pathIndex == len(path) {
				return &t.nodes[i].data, true
			}

			//store current stats
			stackIndex[pathIndex] = i

			//check next path element
			pathIndex++
			i = t.nodes[i].first
		} else {
			i = t.nodes[i].next
			for i == -1 && pathIndex != -1 {
				//go back to previous path
				pathIndex--
				if pathIndex != -1 {
					i = t.nodes[stackIndex[pathIndex]].next
				}
			}
		}
	}
	return nil, false
}
*/
