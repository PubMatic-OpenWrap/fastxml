package fastxml

import (
	"bytes"
)

type ElementTree = tree[XMLToken]

type XMLReader struct {
	in     []byte
	tree   ElementTree
	parser *XMLTokenizer
}

func NewXMLReader(path *xpath) *XMLReader {
	xr := &XMLReader{
		parser: NewXMLTokenizer(path),
	}
	xr.tree = ElementTree{match: xr.match}
	return xr
}

func (xr *XMLReader) match(name string, token XMLToken) bool {
	return bytes.Equal(token.Name(xr.in), []byte(name))
}

func (xr *XMLReader) tokenHandler(name string, parent *Element, child Element) {
	xr.tree.insert(parent, child)
}

func (xr *XMLReader) Parse(in []byte) error {
	xr.tree.reset()
	xr.in = in
	return xr.parser.Parse(in, xr.tokenHandler)
}

func (xr *XMLReader) Childrens(parent *Element) (result []*Element) {
	return xr.tree.getChilds(parent)
}

func (xr *XMLReader) SelectElement(parent *Element, path ...string) *Element {
	if len(path) == 1 {
		return xr.tree.getChild(parent, path[0])
	}
	return xr.tree.getPathNode(parent, path...)
}

func (xr *XMLReader) SelectElements(parent *Element, path ...string) (result []*Element) {
	if len(path) == 1 {
		return xr.tree.getAllChild(parent, path[0])
	}
	return xr.tree.getPathNodes(parent, path...)
}

func (xr *XMLReader) SelectAttr(node *Element, key string) *Attribute {
	attr := node.data.ParseAttribute(xr.in)
	for _, at := range attr {
		if bytes.Equal(at.Key(xr.in), []byte(key)) {
			return &at
		}
	}
	return nil
}

func (xr *XMLReader) SelectAttrValue(node *Element, key string) (value string) {
	if attr := xr.SelectAttr(node, key); attr != nil {
		return string(attr.Value(xr.in))
	}
	return ""
}

func (xr *XMLReader) Text(node *Element, removeCDATA bool) (value string) {
	return string(node.data.Text(xr.in, removeCDATA))
}

func (xr *XMLReader) getXML(in []byte) string {
	buf := bytes.Buffer{}
	start := 0
	for _, node := range xr.tree.nodes {
		buf.Write(in[start:node.data.end.ei])
		start = node.data.end.ei
	}
	return buf.String()
}
