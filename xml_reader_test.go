package fastxml

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMLReader(t *testing.T) {
	in := []byte(`<![CDATA[http://hostedvasttag.url&k=v]]>`)

	xmlReader := NewXMLReader(nil)
	xmlReader.Parse(in[:])

	actual := xmlReader.getXML(in[:])
	t.Logf("Raw Tags: \n%v\n", xmlReader.tree.printRaw(func(t XMLToken) string {
		return fmt.Sprintf("%s:end(%d:%d)", t.Name(in[:]), t.end.si, t.end.ei)
	}))
	t.Logf("XML: %v\n", actual)
	assert.Equal(t, string(in), actual)
}

func TestXMLReader1(t *testing.T) {
	xmldoc := []byte(`<?xml version='1.0'?>
<Catalog>
	<Book id="bk101">
		<Author>Garghentini, Davide</Author>
		<Title>XML Developer's Guide</Title>
		<Genre>Computer</Genre>
		<Price>44.95</Price>
		<PublishDate>2000-10-01</PublishDate>
		<Description>An in-depth look at creating applications with XML.</Description>
	</Book>
	<Book1 id="bk102">
		<Author>Garcia, Debra</Author>
		<Title>Midnight Rain</Title>
		<Genre>Fantasy</Genre>
		<Price>5.95</Price>
		<PublishDate>2000-12-16</PublishDate>
		<Description>A former architect battles corporate zombies, an evil sorceress, and her own childhood to become queen of the world.</Description>
	</Book1>
</Catalog>
	`)

	xmlReader := NewXMLReader(nil)
	err := xmlReader.Parse(xmldoc[:])
	if err != nil {
		t.Errorf("xml parsing error: %s", err.Error())
		return
	}

	t.Logf("\nXML:\n%s", xmldoc)

	for _, element := range xmlReader.SelectElements(nil, "Catalog", "Book") {
		t.Logf("\n/Catalog/Book: id:[%v] innerxml:[%v]", xmlReader.SelectAttrValue(element, "id", ""), xmlReader.Text(element))
	}

	for _, element := range xmlReader.SelectElements(nil, "Catalog", "*", "Author") {
		t.Logf("\n/Catalog/*/Author = %v", xmlReader.Text(element))
	}
}

func TestXMLReader2(t *testing.T) {
	xmldoc := []byte(`<?xml version='1.0'?>
<Catalog>
	<Book id="bk101">
		<Author>Garghentini, Davide</Author>
		<Title>XML Developer's Guide</Title>
		<Genre>Computer</Genre>
		<Price>44.95</Price>
		<PublishDate>2000-10-01</PublishDate>
		<Description>An in-depth look at creating applications with XML.</Description>
	</Book>
	<Book id="bk102">
		<Author>Garcia, Debra</Author>
		<Title>Midnight Rain</Title>
		<Genre>Fantasy</Genre>
		<Price>5.95</Price>
		<PublishDate>2000-12-16</PublishDate>
		<Description>A former architect battles corporate zombies, an evil sorceress, and her own childhood to become queen of the world.</Description>
	</Book>
</Catalog>
	`)

	xmlReader := NewXMLReader(
		GetXPath([][]string{
			{"Catalog", "Book", "Author"},
			{"Catalog", "Book", "Title"},
		}),
	)
	err := xmlReader.Parse(xmldoc[:])
	if err != nil {
		t.Errorf("xml parsing error: %s", err.Error())
		return
	}

	t.Logf("\nXML:\n%s", xmldoc)

	for i, element := range xmlReader.SelectElements(nil, "Catalog", "Book", "Author") {
		t.Logf("\n/Catalog/Book/Author[%d] = %v", i, xmlReader.Text(element))
	}

	for i, element := range xmlReader.SelectElements(nil, "Catalog", "Book", "Genre") {
		t.Logf("\n/Catalog/Book/Genre[%d] = %v", i, xmlReader.Text(element))
	}
}
