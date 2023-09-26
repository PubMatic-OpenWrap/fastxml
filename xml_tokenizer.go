package xmlparser

type TokenNode = treeNode[XMLToken]
type TokenTree = tree[XMLToken]
type TokenHandler func(string, *TokenNode, TokenNode)

func NewTokenNode(token XMLToken) TokenNode {
	return TokenNode{data: token, first: -1, last: -1, next: -1}
}

type XMLTokenizer struct {
	path *xpath
}

func NewXMLTokenizer(path *xpath) *XMLTokenizer {
	return &XMLTokenizer{
		path: path,
	}
}

func (sp *XMLTokenizer) Parse(in []byte, cb TokenHandler) {
	var s stack[TokenNode] //TODOV: get from pool
	var xp stack[*xpath]   //TODOV: get from pool, iff xp.path present

	for i := 0; i < len(in); {
		if in[i] == '<' {
			//get token type
			ttype := getTokenType(in, i+1)

			//get token endindex //TODO this should return token with all details
			endIndex, inlineToken := getTokenEndIndex(in, i, ttype)

			//invalid token
			if endIndex == -1 {
				//there is an issue, append till end and loop will end
				endIndex = len(in)
			}

			if inlineToken {
				ttype = EndXMLToken
			}

			if ttype == StartXMLToken {
				//push start tag into stack and check only for endtags if those are matching to ours tag
				token := XMLToken{
					start: xmlTagIndex{si: i, ei: endIndex},
				}

				//xpath handling
				if sp.path != nil && s.len() == xp.len() {
					path := xp.peek()
					if path == nil {
						path = &sp.path
					}

					/*NOTE: do not use existing path, it will update stack variable*/
					p := (*path).get(string(token.Name(in)))
					if p != nil {
						xp.push(p)
					}
				}

				s.push(TokenNode{data: token, first: -1, last: -1, next: -1})
			} else if ttype == EndXMLToken {
				//get start xml tag
				foundTag := true
				child := s.pop()
				child.data.end = xmlTagIndex{si: i, ei: endIndex}
				//fmt.Printf("%s:<%d,%d,%d>\n", string(child.data.Name(in)), child.data.start.si, child.data.end.ei, child.data.end.ei-child.data.start.si)

				//xpath handling
				if sp.path != nil {
					if s.len() < xp.len() {
						xp.pop()
					} else {
						foundTag = false
					}
				}

				if foundTag && cb != nil {
					//append tokens to list
					cb(string(child.data.Name(in[:])), s.peek(), child)
				}
			}
			i = endIndex
			continue
		}
		i++
	}
}
