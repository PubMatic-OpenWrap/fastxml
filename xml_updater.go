package fastxml

import (
	"sort"
)

type xmlOperation struct {
	si, ei int
	data   XMLWriter
}

/*
TODO: makesure not 2 operations overlaps
*/
type XMLUpdater struct {
	xmlReader     *XMLReader
	writeSettings WriteSettings
	ops           []xmlOperation
}

type WriteSettings struct {
	CDATAWrap    bool
	ExpandInline bool
}

func NewXMLUpdater(xmlReader *XMLReader, writeSettings WriteSettings) *XMLUpdater {
	return &XMLUpdater{xmlReader: xmlReader, writeSettings: writeSettings}
}

/* XML ELEMENT FUNCTION */

/* XML ELEMENT FUNCTION */

func (xu *XMLUpdater) AppendElement(element *Element, tagXML XMLWriter) {
	if element == nil || tagXML == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si:   element.data.end.si,
		ei:   element.data.end.si,
		data: tagXML,
	})
}

func (xu *XMLUpdater) BeforeElement(element *Element, tagXML XMLWriter) {
	if element == nil || tagXML == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si:   element.data.start.si,
		ei:   element.data.start.si,
		data: tagXML,
	})
	/*
		//INLINE TAG NOT SUPPORTED YET
		if element.data.IsInline() {
			it should replace "/>" value with ">tagXML</xmlns:name>"
			need to check if multiple such operation are there then only append tagXML
		}
	*/
}

func (xu *XMLUpdater) AfterElement(element *Element, tagXML XMLWriter) {
	if element == nil || tagXML == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si:   element.data.end.ei,
		ei:   element.data.end.ei,
		data: tagXML,
	})
	/*
		//INLINE TAG NOT SUPPORTED YET
		if element.data.IsInline() {
			it should replace "/>" value with ">tagXML</xmlns:name>"
			need to check if multiple such operation are there then only append tagXML
		}
	*/
}

func (xu *XMLUpdater) PrependElement(element *Element, tagXML XMLWriter) {
	if element == nil || tagXML == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si:   element.data.start.ei,
		ei:   element.data.start.ei,
		data: tagXML,
	})
	/*
		//INLINE TAG NOT SUPPORTED YET
		if element.data.IsInline() {
			it should replace "/>" value with ">tagXML</xmlns:name>"
			need to check if multiple such operation are there then only append tagXML
		}
	*/
}

func (xu *XMLUpdater) ReplaceElement(element *Element, tagXML XMLWriter) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si:   element.data.start.si,
		ei:   element.data.end.ei,
		data: tagXML,
	})
}

func (xu *XMLUpdater) RemoveElement(element *Element) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: element.data.start.si,
		ei: element.data.end.ei,
	})
}

func (xu *XMLUpdater) UpdateText(element *Element, text string, cdata bool, escaping XMLEscapingMode) {
	if element == nil {
		return
	}

	op := xmlOperation{
		si: element.data.start.ei,
		ei: element.data.end.si,
		data: &XMLTextElement{
			text:     []byte(text),
			cdata:    cdata,
			escaping: escaping,
		},
	}

	if element.data.IsInline() {
		if xu.writeSettings.ExpandInline {
			//TODO: BUG DONOT ALLOW UPDATING TEXT INCASE OF ExpandInline TRUE
			return
		}

		op.si = element.data.end.ei - 2
		op.ei = element.data.end.ei - 1
		op.data = NewXmlTextFunc(
			false,
			func(w Writer, args ...any) {
				if len(args) != 2 {
					return
				}
				name, _ := args[0].(string)
				text, _ := args[1].(XMLWriter)

				w.WriteByte('>')
				text.Write(w)
				w.WriteString("</")
				w.WriteString(name)
			},
			xu.xmlReader.NSName(element), op.data)
	}

	xu.ops = append(xu.ops, op)
}

/* ATTRIBUTE FUNCTIONS */
func (xu *XMLUpdater) AddAttribute(element *Element, namespace, key, value string) {
	if element == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: element.data.name.ei,
		ei: element.data.name.ei,
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
		si: attr.key.si - 1,
		ei: attr.value.ei + 1,
	})
}

func (xu *XMLUpdater) UpdateAttributeValue(attr *Attribute, value string) {
	if attr == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: attr.value.si - 1,
		ei: attr.value.ei + 1,
		data: NewXmlTextFunc(
			false,
			func(w Writer, args ...any) {
				if len(args) != 1 {
					return
				}
				value, _ := args[0].(string)
				w.WriteByte('"')
				quoteEscape(w, value)
				w.WriteByte('"')
			},
			value,
		),
	})
}

func (xu *XMLUpdater) expandInline(element *Element) {
	if element == nil || !element.data.IsInline() {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: element.data.end.ei - 2,
		ei: element.data.end.ei - 1,
		data: NewXmlTextFunc(
			false,
			func(w Writer, args ...any) {
				if len(args) != 1 {
					return
				}
				value, _ := args[0].(string)
				w.WriteString("></")
				w.WriteString(value)
			},
			xu.xmlReader.NSName(element),
		),
	})
}

/* //TODO: NOT NEEDED FUNCTION
func (xu *XMLUpdater) UpdateAttributeName(attr *Attribute, key string) {
	if attr == nil {
		return
	}
	xu.ops = append(xu.ops, xmlOperation{
		si: attr.key.si,
		ei: attr.key.ei,
		data: &XMLTextElement{
			text: []byte(key),
		},
	})
}
*/

func (xu *XMLUpdater) ApplyXMLSettingsOperations() {
	// wrap cdata
	if !(xu.writeSettings.ExpandInline || xu.writeSettings.CDATAWrap) {
		return
	}

	xu.xmlReader.Iterate(func(element *Element) {
		if xu.xmlReader.IsLeaf(element) {
			name := xu.xmlReader.Name(element)
			_ = name
			text := xu.xmlReader.RawText(element)
			trimmedText := trimSpace(text) //strings.TrimSpace(text)

			if xu.xmlReader.IsCDATA(element) {
				//wrap text into cdata text => <![CDATA[text]]> or remove only whitespaces
				xu.UpdateText(element, trimmedText, len(trimmedText) != 0, NoEscaping)
			} else {
				if xu.writeSettings.ExpandInline && element.Data().IsInline() {
					//expand inline tag <abc/> => <abc></abc>
					xu.expandInline(element)
					return
				}

				if xu.writeSettings.CDATAWrap {
					//wrap text into cdata text => <![CDATA[text]]> or remove only whitespaces
					xu.UpdateText(element, trimmedText, len(trimmedText) != 0, XMLUnescapeMode)
				}
			}
		}
	})
}

func (xu *XMLUpdater) Build(buf Writer) {
	//sort operations based on index
	sort.SliceStable(xu.ops[:], func(i, j int) bool {
		return (xu.ops[i].si < xu.ops[j].si ||
			(xu.ops[i].si == xu.ops[j].si && xu.ops[i].ei < xu.ops[j].ei))
	})

	in := xu.xmlReader.RawXML()
	offset := 0
	for _, op := range xu.ops {
		if offset <= op.si {
			buf.Write(in[offset:op.si])
			offset = op.ei
		}
		if op.data != nil {
			op.data.Write(buf)
		}
	}
	buf.Write(in[offset:])
}
