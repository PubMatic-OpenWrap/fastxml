package fastxml

import (
	"bytes"
	"fmt"
)

type xmlWriter interface {
	Write(buf *bytes.Buffer)
}

// xmlText element
type xmlText struct {
	text []byte
}

func (xt *xmlText) Write(buf *bytes.Buffer) {
	buf.Write(xt.text)
}

// xmlAttribute element
type xmlAttribute struct {
	namespace, key, value string
}

func (xa *xmlAttribute) Write(buf *bytes.Buffer) {
	buf.WriteByte(' ')
	if xa.namespace != "" {
		buf.WriteString(xa.namespace)
		buf.WriteByte(':')
	}
	buf.WriteString(xa.key)
	buf.WriteString(`="`)
	buf.WriteString(xa.value)
	buf.WriteByte('"')
}

// xmlTag element
type xmlTag struct {
	attr  []xmlAttribute
	name  string
	text  string
	child []*xmlTag
}

func NewXMLTag(name string, text string) *xmlTag {
	return &xmlTag{name: name, text: text}
}

func (xt *xmlTag) AddAttribute(namespace, key, value string) *xmlTag {
	xt.attr = append(xt.attr, xmlAttribute{namespace: namespace, key: key, value: value})
	return xt
}

func (xt *xmlTag) AddChild(child *xmlTag) *xmlTag {
	xt.child = append(xt.child, child)
	return xt
}

func (xt *xmlTag) Write(buf *bytes.Buffer) {
	if len(xt.name) > 0 {
		buf.WriteByte('<')
		buf.WriteString(xt.name)
		for _, attr := range xt.attr {
			attr.Write(buf)
		}
		buf.WriteByte('>')
	}

	if len(xt.text) > 0 {
		buf.WriteString(xt.text)
	}

	for _, child := range xt.child {
		child.Write(buf)
	}

	if len(xt.name) > 0 {
		fmt.Fprintf(buf, `</%s>`, xt.name)
	}
}
