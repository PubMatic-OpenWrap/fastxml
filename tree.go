package fastxml

import (
	"bytes"
	"fmt"
	"strings"
)

type compareXMLToken func(string, XMLToken) bool

type treeNode struct {
	data                   XMLToken
	idx, first, last, next int
}

func (n treeNode) Data() XMLToken {
	return n.data
}

func (n treeNode) Index() int {
	return n.idx
}

func (n treeNode) IsLeaf() bool {
	return n.first == -1
}

type xmlTree struct {
	nodes []treeNode //first node will be always last node
	match compareXMLToken
}

/*
insert function to insert node n in parent node
NOTE: always re-fetch parent object everytime when inserting new object
*/
func (t *xmlTree) insert(parent *treeNode, n treeNode) {
	if len(t.nodes) == 0 {
		t.nodes = append(t.nodes, treeNode{idx: 0, first: -1, last: -1, next: -1})
	}

	n.idx = len(t.nodes)
	if parent == nil {
		parent = &t.nodes[0]
	}

	if parent.first == -1 {
		//first node
		parent.first = n.idx
	} else {
		//subsequent node
		t.nodes[parent.last].next = n.idx
	}
	parent.last = n.idx

	t.nodes = append(t.nodes, n)
}

func (t *xmlTree) reset() {
	t.nodes = t.nodes[:0]
}

func (t *xmlTree) getChild(parent *treeNode, child string) (result *treeNode) {
	parentIndex := 0
	if parent != nil {
		parentIndex = parent.idx
	}

	if parentIndex >= len(t.nodes) {
		return nil
	}

	for i := t.nodes[parentIndex].first; i != -1; i = t.nodes[i].next {
		if t.match != nil && t.match(child, t.nodes[i].data) {
			return &t.nodes[i]
		}
	}
	return nil
}

func (t *xmlTree) getAllChild(parent *treeNode, child string) (result []*treeNode) {
	parentIndex := 0
	if parent != nil {
		parentIndex = parent.idx
	}

	if parentIndex >= len(t.nodes) {
		return nil
	}

	for i := t.nodes[parentIndex].first; i != -1; i = t.nodes[i].next {
		if t.match != nil && t.match(child, t.nodes[i].data) {
			result = append(result, &t.nodes[i])
		}
	}
	return
}

func (t *xmlTree) getChilds(parent *treeNode) (result []*treeNode) {
	parentIndex := 0
	if parent != nil {
		parentIndex = parent.idx
	}

	if parentIndex >= len(t.nodes) {
		return nil
	}

	for i := t.nodes[parentIndex].first; i != -1; i = t.nodes[i].next {
		result = append(result, &t.nodes[i])
	}
	return
}

func (t *xmlTree) _getPath(parent int, result *[]*treeNode, path ...string) {
	for i := t.nodes[parent].first; i != -1; i = t.nodes[i].next {
		if t.match != nil && t.match(path[0], t.nodes[i].data) {
			if len(path) == 1 {
				(*result) = append((*result), &t.nodes[i])
			} else {
				t._getPath(i, result, path[1:]...)
			}
		}
	}
}

func (t *xmlTree) getPathNodes(parent *treeNode, path ...string) (result []*treeNode) {
	parentIndex := 0
	if parent != nil {
		parentIndex = parent.idx
	}
	if parentIndex >= len(t.nodes) {
		return nil
	}
	t._getPath(parentIndex, &result, path...)
	return
}

func (t *xmlTree) getPathNode(parent *treeNode, path ...string) (result *treeNode) {
	parentIndex := 0
	if parent != nil {
		parentIndex = parent.idx
	}

	if parentIndex >= len(t.nodes) {
		return nil
	}

	for iPath := 0; iPath < len(path); iPath++ {
		j := t.nodes[parentIndex].first
		for j != -1 {
			if t.match != nil && t.match(path[iPath], t.nodes[j].data) {
				//found
				break
			}
			j = t.nodes[j].next
		}
		if j == -1 {
			//not found
			return nil
		}
		parentIndex = j
	}

	return &t.nodes[parentIndex]
}

func (t *xmlTree) iterate(f func(*treeNode)) {
	for i := range t.nodes {
		f(&t.nodes[i])
	}
}

func (t *xmlTree) _traverse(index int, f func(*treeNode)) {
	f(&t.nodes[index])
	for i := t.nodes[index].first; i != -1; i = t.nodes[i].next {
		t._traverse(i, f)
	}
}

func (t *xmlTree) traverse(node *treeNode, f func(*treeNode)) {
	parent := 0
	if node != nil {
		parent = node.idx
	}
	t._traverse(parent, f)
}

/* Printing Function */
func (t *xmlTree) _print(buf *bytes.Buffer, index, indent int, f func(XMLToken) string) {
	buf.WriteByte('\n')
	buf.WriteString(strings.Repeat("\t", indent))
	buf.WriteByte('|')
	buf.WriteString(f(t.nodes[index].data))
	indent++
	for i := t.nodes[index].first; i != -1; i = t.nodes[i].next {
		t._print(buf, i, indent, f)
	}
}

func (t *xmlTree) print(f func(XMLToken) string) string {
	root := len(t.nodes) - 1
	buf := bytes.Buffer{}
	t._print(&buf, root, 0, f)
	return buf.String()
}

func (t *xmlTree) printRaw(f func(XMLToken) string) string {
	buf := bytes.Buffer{}
	for i, node := range t.nodes {
		buf.WriteString(fmt.Sprintf("\n%d:<%d,%d,%d>:%s", i, node.first, node.last, node.next, f(node.data)))
	}
	return buf.String()
}
