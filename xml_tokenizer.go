package fastxml

import (
	"fmt"
)

var errInvalidXML = fmt.Errorf("invalid xml")

type Element = treeNode[XMLToken]
type TokenHandler func(string, *Element, Element)

type XMLTokenizer struct{}

func NewXMLTokenizer() *XMLTokenizer {
	return &XMLTokenizer{}
}

func (sp *XMLTokenizer) parse(in []byte, cb TokenHandler) error {
	/*
		TODO:
		1. get s from pool,
		2. get xp from pool, iff xp.path present
		3. adding defer to put stack and xp into pool back
	*/
	var s stack[Element]

	for i := 0; i < len(in); {
		if in[i] == '<' {
			//get token type
			ttype := getTokenType(in, i+1)

			//TODO this should return token with all details
			//get token endindex
			endIndex, inlineToken := getTokenEndIndex(in, i+1, ttype)

			//invalid token
			if endIndex == -1 {
				return errInvalidXML
			}

			if inlineToken {
				ttype = endXMLToken
			}

			if ttype == startXMLToken {
				//push start tag into stack and check only for endtags if those are matching to ours tag
				token := XMLToken{
					start: xmlTagIndex{si: i, ei: endIndex},
				}
				s.push(Element{data: token, first: -1, last: -1, next: -1})
			} else if ttype == endXMLToken {
				//get start xml tag
				var child *Element

				if inlineToken {
					child = &Element{
						data: XMLToken{
							start: xmlTagIndex{si: i, ei: endIndex},
						},
						first: -1, last: -1, next: -1,
					}
				} else {
					child = s.pop()
				}

				if child == nil {
					return errInvalidXML
				}
				child.data.end = xmlTagIndex{si: i, ei: endIndex}

				if cb != nil {
					//append tokens to list
					cb(string(child.data.Name(in[:])), s.peek(), *child)
				}
				//fmt.Printf("%s:<%d,%d,%d>\n", string(child.data.Name(in)), child.data.start.si, child.data.end.ei, child.data.end.ei-child.data.start.si)
			}
			i = endIndex
			continue
		}
		i++
	}
	if s.len() != 0 {
		return errInvalidXML
	}
	return nil
}

func (sp *XMLTokenizer) parseUsingXPath(in []byte, ixpath *xpath, cb TokenHandler) error {
	/*
		TODO:
		1. get s from pool,
		2. get xp from pool, iff xp.path present
		3. adding defer to put stack and xp into pool back
	*/
	var s stack[Element]
	var xp stack[*xpath]

	for i := 0; i < len(in); {
		if in[i] == '<' {
			//get token type
			ttype := getTokenType(in, i+1)

			//TODO this should return token with all details
			//get token endindex
			endIndex, inlineToken := getTokenEndIndex(in, i+1, ttype)

			//invalid token
			if endIndex == -1 {
				return errInvalidXML
			}

			if inlineToken {
				ttype = endXMLToken
			}

			if ttype == startXMLToken {
				//push start tag into stack and check only for endtags if those are matching to ours tag
				token := XMLToken{
					start: xmlTagIndex{si: i, ei: endIndex},
				}

				//xpath handling
				if ixpath != nil && s.len() == xp.len() {
					path := xp.peek()
					if path == nil {
						path = &ixpath
					}

					/*NOTE: do not use existing path, it will update stack variable*/
					p := (*path).get(string(token.Name(in)))
					if p != nil {
						xp.push(p)
					}
				}

				s.push(Element{data: token, first: -1, last: -1, next: -1})
			} else if ttype == endXMLToken {
				//get start xml tag
				foundTag := true
				var child *Element

				if inlineToken {
					child = &Element{
						data: XMLToken{
							start: xmlTagIndex{si: i, ei: endIndex},
						},
						first: -1, last: -1, next: -1,
					}
				} else {
					child = s.pop()
				}

				if child == nil {
					return errInvalidXML
				}
				child.data.end = xmlTagIndex{si: i, ei: endIndex}

				//xpath handling
				if ixpath != nil {
					if s.len() < xp.len() {
						xp.pop()
					} else {
						foundTag = false
					}
				}

				if foundTag && cb != nil {
					//append tokens to list
					cb(string(child.data.Name(in[:])), s.peek(), *child)
				}
				//fmt.Printf("%s:<%d,%d,%d>\n", string(child.data.Name(in)), child.data.start.si, child.data.end.ei, child.data.end.ei-child.data.start.si)
			}
			i = endIndex
			continue
		}
		i++
	}
	if s.len() != 0 {
		return errInvalidXML
	}
	return nil
}

/*
func NewElement(token XMLToken) Element {
	return Element{data: token, first: -1, last: -1, next: -1}
}
*/
