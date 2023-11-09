package fastxml

import (
	"bytes"
	"sort"
)

type xmlOperationType int

const (
	defaultXMLOperationType xmlOperationType = iota
	addAttributeXMLOperationType
)

type xmlOperation struct {
	op         xmlOperationType
	si, ei     int
	parameters []any //TODO: replace with string variable
}

/*
TODO: makesure not 2 operations overlaps
*/
type XMLUpdater struct {
	in  []byte
	ops []xmlOperation
}

func NewXMLUpdater(in []byte) *XMLUpdater {
	return &XMLUpdater{in: in[:]}
}

/* XML ELEMENT FUNCTION */

func (xu *XMLUpdater) AppendElement(element *Element, tagXML string) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{si: element.Data().end.si, ei: element.Data().end.si, parameters: []any{tagXML}})
}

func (xu *XMLUpdater) PrependElement(element *Element, tagXML string) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{si: element.Data().start.ei, ei: element.Data().start.ei, parameters: []any{tagXML}})
	/*
		//INLINE TAG NOT SUPPORTED YEY
		if element.Data().IsInline() {
			it should replace "/>" value with ">tagXML</xmlns:name>"
			need to check if multiple such operation are there then only append tagXML
		}
	*/
}

func (xu *XMLUpdater) ReplaceElement(element *Element, tagXML string) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{si: element.Data().start.si, ei: element.Data().end.ei, parameters: []any{tagXML}})
}

func (xu *XMLUpdater) RemoveElement(element *Element) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{si: element.Data().start.si, ei: element.data.end.ei})
}

func (xu *XMLUpdater) UpdateText(element *Element, text string) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{si: element.Data().start.ei, ei: element.Data().end.si, parameters: []any{text}})
}

/* ATTRIBUTE FUNCTIONS */
func (xu *XMLUpdater) AddAttribute(element *Element, namespace, key, value string) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		op:         addAttributeXMLOperationType,
		si:         element.Data().name.ei,
		ei:         element.Data().name.ei,
		parameters: []any{namespace, key, value}})
}

func (xu *XMLUpdater) RemoveAttribute(attr *Attribute) {
	if attr == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{si: attr.key.si, ei: attr.value.ei + 1})
}

func (xu *XMLUpdater) UpdateAttributeName(attr *Attribute, key string) {
	if attr == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{si: attr.key.si, ei: attr.key.ei, parameters: []any{key}})
}

func (xu *XMLUpdater) UpdateAttributeValue(attr *Attribute, value string) {
	if attr == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{si: attr.value.si, ei: attr.value.ei, parameters: []any{value}})
}

func (xu *XMLUpdater) update(buf *bytes.Buffer, xop xmlOperation) {
	switch xop.op {
	case addAttributeXMLOperationType:
		ns := xop.parameters[0].(string)
		key := xop.parameters[1].(string)
		value := xop.parameters[2].(string)

		buf.WriteByte(' ') //writing space before adding attribute
		if ns != "" {
			buf.WriteString(ns)
			buf.WriteByte(':')
		}
		buf.WriteString(key)
		buf.WriteString(`="`)
		buf.WriteString(value)
		buf.WriteByte('"')
	default:
		if len(xop.parameters) > 0 {
			for _, param := range xop.parameters {
				buf.WriteString(param.(string))
			}
		}
	}
}

func (xu *XMLUpdater) Build(buf *bytes.Buffer) {
	//sort operations based on index
	sort.SliceStable(xu.ops[:], func(i, j int) bool {
		return (xu.ops[i].si < xu.ops[j].si ||
			(xu.ops[i].si == xu.ops[j].si && xu.ops[i].ei < xu.ops[j].ei))
	})

	offset := 0
	for _, op := range xu.ops {
		if offset <= op.si {
			buf.Write(xu.in[offset:op.si])
			offset = op.ei
		}
		if len(op.parameters) > 0 {
			xu.update(buf, op)
		}
	}
	buf.Write(xu.in[offset:])
}
