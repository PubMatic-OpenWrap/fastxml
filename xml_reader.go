package xmlparser

import "bytes"

type XMLReader struct {
	tree   TokenTree
	parser *XMLTokenizer
}

func NewXMLReader(path *xpath) *XMLReader {
	return &XMLReader{
		tree:   TokenTree{},
		parser: NewXMLTokenizer(path),
	}
}

func (xr *XMLReader) tokenHandler(name string, parent *TokenNode, child TokenNode) {
	xr.tree.insertChild(parent, child)
}

func (xr *XMLReader) Parse(in []byte) {
	xr.tree.reset()
	xr.parser.Parse(in, xr.tokenHandler)
}

func (xr *XMLReader) GetXML(in []byte) string {
	buf := bytes.Buffer{}
	start := 0
	for _, node := range xr.tree.nodes {
		buf.Write(in[start:node.data.end.ei])
		start = node.data.end.ei
	}
	return buf.String()
}
