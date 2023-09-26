package fastxml

import (
	"bytes"
	"sort"
)

type xmlOperation struct {
	si, ei int
	data   xmlWriter
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

func (xu *XMLUpdater) AppendElement(element *Element, tagXML *xmlTag) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si:   element.Data().end.si,
		ei:   element.Data().end.si,
		data: tagXML,
	})
}

func (xu *XMLUpdater) PrependElement(element *Element, tagXML *xmlTag) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si:   element.Data().start.ei,
		ei:   element.Data().start.ei,
		data: tagXML,
	})
	/*
		//INLINE TAG NOT SUPPORTED YEY
		if element.Data().IsInline() {
			it should replace "/>" value with ">tagXML</xmlns:name>"
			need to check if multiple such operation are there then only append tagXML
		}
	*/
}

func (xu *XMLUpdater) ReplaceElement(element *Element, tagXML *xmlTag) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si:   element.Data().start.si,
		ei:   element.Data().end.ei,
		data: tagXML,
	})
}

func (xu *XMLUpdater) RemoveElement(element *Element) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: element.Data().start.si,
		ei: element.data.end.ei,
	})
}

func (xu *XMLUpdater) UpdateText(element *Element, text string) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: element.Data().start.ei,
		ei: element.Data().end.si,
		data: &xmlText{
			text: []byte(text),
		},
	})
}

/* ATTRIBUTE FUNCTIONS */
func (xu *XMLUpdater) AddAttribute(element *Element, namespace, key, value string) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: element.Data().name.ei,
		ei: element.Data().name.ei,
		data: &xmlAttribute{
			namespace: namespace,
			key:       key,
			value:     value,
		},
	})
}

func (xu *XMLUpdater) RemoveAttribute(attr *Attribute) {
	if attr == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: attr.key.si,
		ei: attr.value.ei + 1,
	})
}

func (xu *XMLUpdater) UpdateAttributeName(attr *Attribute, key string) {
	if attr == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: attr.key.si,
		ei: attr.key.ei,
		data: &xmlText{
			text: []byte(key),
		},
	})
}

func (xu *XMLUpdater) UpdateAttributeValue(attr *Attribute, value string) {
	if attr == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: attr.value.si,
		ei: attr.value.ei,
		data: &xmlText{
			text: []byte(value),
		},
	})
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
		if op.data != nil {
			op.data.Write(buf)
		}
	}
	buf.Write(xu.in[offset:])
}
