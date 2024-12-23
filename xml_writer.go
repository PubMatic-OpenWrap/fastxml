package fastxml

import (
	"io"
)

type Writer interface {
	io.StringWriter
	io.ByteWriter
	io.Writer
}

type XMLWriter interface {
	Write(buf Writer)
}

//---------------------------------------------------------------------------------------------

// xmlAttribute element
type xmlAttribute struct {
	namespace, key, value string
}

func (xa *xmlAttribute) Write(buf Writer) {
	buf.WriteByte(' ')
	if xa.namespace != "" {
		buf.WriteString(xa.namespace)
		buf.WriteByte(':')
	}
	buf.WriteString(xa.key)
	buf.WriteString(`="`)
	quoteEscape(buf, xa.value)
	buf.WriteByte('"')
}

//---------------------------------------------------------------------------------------------

// XMLElement element
type XMLElement struct {
	ns    string
	name  string
	text  XMLWriter
	attr  []xmlAttribute
	child []XMLWriter
}

func CreateElement(name string) *XMLElement {
	return &XMLElement{name: name}
}

func (xt *XMLElement) SetNamespace(namespace string) *XMLElement {
	xt.ns = namespace
	return xt
}

func (xt *XMLElement) AddAttribute(namespace, key, value string) *XMLElement {
	xt.attr = append(xt.attr, xmlAttribute{namespace: namespace, key: key, value: value})
	return xt
}

func (xt *XMLElement) AddChild(child XMLWriter) *XMLElement {
	xt.child = append(xt.child, child)
	return xt
}

func (xt *XMLElement) SetText(text string, cdata bool, escaping XMLEscapingMode) *XMLElement {
	xt.text = &XMLTextElement{text: []byte(text), cdata: cdata, escaping: escaping}
	return xt
}

func (xt *XMLElement) SetName(name string) *XMLElement {
	xt.name = name
	return xt
}

func (xt *XMLElement) Write(buf Writer) {
	if len(xt.name) > 0 {
		buf.WriteByte('<')

		if len(xt.ns) > 0 {
			buf.WriteString(xt.ns)
			buf.WriteByte(':')
		}

		buf.WriteString(xt.name)

		for _, attr := range xt.attr {
			attr.Write(buf)
		}

		buf.WriteByte('>')
	}

	if xt.text != nil {
		xt.text.Write(buf)
	}

	for _, child := range xt.child {
		child.Write(buf)
	}

	if len(xt.name) > 0 {
		buf.WriteString("</")
		if len(xt.ns) > 0 {
			buf.WriteString(xt.ns)
			buf.WriteByte(':')
		}
		buf.WriteString(xt.name)
		buf.WriteByte('>')
	}
}

//---------------------------------------------------------------------------------------------

// XMLTextElement element
type XMLTextElement struct {
	cdata    bool
	escaping XMLEscapingMode
	text     []byte
}

func NewXMLText(text string, cdata bool, escaping XMLEscapingMode) *XMLTextElement {
	return &XMLTextElement{
		text:     []byte(text),
		cdata:    cdata,
		escaping: escaping,
	}
}

func NewXMLBytes(text []byte, cdata bool, escaping XMLEscapingMode) *XMLTextElement {
	return &XMLTextElement{
		text:     text,
		cdata:    cdata,
		escaping: escaping,
	}
}

func (xt *XMLTextElement) Write(buf Writer) {
	if xt.cdata {
		buf.Write(cdataStart)
	}

	// Write Text
	switch xt.escaping {
	case XMLEscapeMode:
		escape(buf, xt.text)
	case XMLUnescapeMode:
		unescape(buf, xt.text)
	default:
		buf.Write(xt.text)
	}

	if xt.cdata {
		buf.Write(cdataEnd)
	}
}

//---------------------------------------------------------------------------------------------

// XMLTextFunc element
type XMLTextFunc struct {
	cdata bool
	fn    func(Writer, ...any)
	args  []interface{}
}

func NewXmlTextFunc(cdata bool, f func(Writer, ...any), args ...any) *XMLTextFunc {
	return &XMLTextFunc{
		cdata: cdata,
		fn:    f,
		args:  args,
	}
}

func (xf *XMLTextFunc) Write(buf Writer) {
	if xf.fn == nil {
		return
	}

	if xf.cdata {
		buf.Write(cdataStart)
	}
	xf.fn(buf, xf.args...)
	if xf.cdata {
		buf.Write(cdataEnd)
	}
}
