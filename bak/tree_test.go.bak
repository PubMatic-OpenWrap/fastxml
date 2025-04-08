package fastxml

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_insert(t *testing.T) {
	iTree := &tree[int]{
		match: func(str string, value int) bool {
			return str == fmt.Sprint(value)
		},
	}

	emptyTreeNode := func(data int) treeNode[int] {
		return treeNode[int]{
			data:  data,
			index: 0,
			first: -1,
			last:  -1,
			next:  -1,
		}
	}

	nodes2array := func(nodes []*treeNode[int]) (values []int) {
		for i := 0; i < len(nodes); i++ {
			values = append(values, nodes[i].Data())
		}
		return
	}

	assert.Equal(t, 0, len(iTree.nodes))

	//adding 1st level elements
	iTree.insert(nil, emptyTreeNode(10))
	iTree.insert(nil, emptyTreeNode(20))
	iTree.insert(nil, emptyTreeNode(30))

	//adding 2nd level elements
	iTree.insert(&iTree.nodes[1], emptyTreeNode(11))
	iTree.insert(&iTree.nodes[1], emptyTreeNode(12))
	iTree.insert(&iTree.nodes[2], emptyTreeNode(21))
	iTree.insert(&iTree.nodes[2], emptyTreeNode(22))
	iTree.insert(&iTree.nodes[3], emptyTreeNode(300))
	iTree.insert(&iTree.nodes[3], emptyTreeNode(301))
	iTree.insert(&iTree.nodes[3], emptyTreeNode(300))

	//assert tree
	childs := iTree.getChilds(nil)
	assert.Equal(t, []int{10, 20, 30}, nodes2array(childs))

	child := iTree.getChild(nil, "10")
	assert.NotNil(t, child)
	assert.Equal(t, child.Data(), 10)
	assert.Equal(t, []int{11, 12}, nodes2array(iTree.getChilds(child)))

	child = iTree.getChild(nil, "20")
	assert.NotNil(t, child)
	assert.Equal(t, child.Data(), 20)
	assert.Equal(t, []int{21, 22}, nodes2array(iTree.getChilds(child)))

	child = iTree.getChild(nil, "30")
	assert.NotNil(t, child)
	assert.Equal(t, child.Data(), 30)
	assert.Equal(t, []int{300, 301, 300}, nodes2array(iTree.getChilds(child)))

	childs = iTree.getAllChild(iTree.getChild(nil, "30"), "300")
	assert.Equal(t, []int{300, 300}, nodes2array(childs))

	assert.Equal(t, 11, len(iTree.nodes))

	//getPathNode
	child = iTree.getPathNode(nil, "20", "21")
	assert.NotNil(t, child)
	assert.Equal(t, child.data, 21)

	child = iTree.getPathNode(nil, "30", "301")
	assert.NotNil(t, child)
	assert.Equal(t, child.data, 301)

	//not found
	child = iTree.getPathNode(nil, "30", "302")
	assert.Nil(t, child)

	//same multiple nodes
	childs = iTree.getPathNodes(nil, "30", "300")
	assert.Equal(t, []int{300, 300}, nodes2array(childs))

	//not found
	childs = iTree.getPathNodes(nil, "30", "302")
	assert.Equal(t, len(childs), 0)

	iTree.reset()
	assert.Equal(t, 0, len(iTree.nodes))
}
