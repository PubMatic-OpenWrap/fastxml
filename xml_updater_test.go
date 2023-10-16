package fastxml

import (
	"bytes"
	"fmt"
	"testing"
)

func TestXMLUpdater(t *testing.T) {
	xmlDoc := []byte(`
<a>
    <b>b-data</b>
    <c>c-data</c>
    <d>
        <e>e-data</e>
		<c>c-data</c>
    </d>
	<f>
		<g>g-data</g>
	</f>
</a>`)

	reader := NewXMLReader(nil)
	if err := reader.Parse(xmlDoc); err != nil {
		fmt.Printf("xml parsing error: %v", err.Error())
		return
	}

	elementB := reader.FindElement(nil, "a", "b")
	elementG := reader.FindElement(nil, "a", "f", "g")

	//xmlUpdater
	updater := NewXMLUpdater(
		xmlDoc,
		func(buf *bytes.Buffer, args ...any) {
			//CONTROL IS WITH USER TO DO CUSTOM MODIFICATION
			buf.WriteString(args[0].(string))
		},
	)

	/*
		Replace function to replace any thing from start to end index
		it can be tag, attribute, attribute value, text, xml
	*/
	start, end := elementB.Data().TagOffset()
	updater.Replace(start, end, "<new_b>new_b_data</new_b>")

	/*
		Insert function to insert tag at specific index
	*/
	start, end = elementG.Data().TagOffset()
	updater.Insert(start, "<g1>prepend_data</g1>")
	updater.Insert(end, "<g2>append_data</g2>")

	//Build Updated XML File
	buf := bytes.Buffer{}
	updater.Build(&buf)

	fmt.Printf("\nOriginal XML:%s", xmlDoc)
	fmt.Printf("\n\nUpdated XML:%s", buf.String())
}
