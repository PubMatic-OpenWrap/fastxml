****** REPOSITORY UNDER DEVELOPMENT ******

# xmlparser
string based xmlparser utility

# Usages
## Reading XML File
```
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

	xmlReader := NewXMLReader(nil)
	err := xmlReader.Parse(xmldoc[:])
	if err != nil {
		fmt.Printf("xml parsing error: %s", err.Error())
		return
	}

	fmt.Printf("\nXML:\n%s", xmldoc)

	for _, element := range xmlReader.FindElements(nil, "Catalog", "Book") {
		fmt.Printf("\n/Catalog/Book: id:[%v] innerxml:[%v]", xmlReader.GetAttribute(element, "id"), xmlReader.GetText(element, true))
	}

	for _, element := range xmlReader.FindElements(nil, "Catalog", "Book", "Author") {
		fmt.Printf("\n/Catalog/Book/Author = %v", xmlReader.GetText(element, true))
	}

Output:
XML:
<?xml version='1.0'?>
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
	
/Catalog/Book: id:[bk101] innerxml:[
		<Author>Garghentini, Davide</Author>
		<Title>XML Developer's Guide</Title>
		<Genre>Computer</Genre>
		<Price>44.95</Price>
		<PublishDate>2000-10-01</PublishDate>
		<Description>An in-depth look at creating applications with XML.</Description>
	]
/Catalog/Book: id:[bk102] innerxml:[
		<Author>Garcia, Debra</Author>
		<Title>Midnight Rain</Title>
		<Genre>Fantasy</Genre>
		<Price>5.95</Price>
		<PublishDate>2000-12-16</PublishDate>
		<Description>A former architect battles corporate zombies, an evil sorceress, and her own childhood to become queen of the world.</Description>
	]
/Catalog/Book/Author = Garghentini, Davide
/Catalog/Book/Author = Garcia, Debra
```

## Reading Specific Tags Only
```
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
		fmt.Printf("xml parsing error: %s", err.Error())
		return
	}

	fmt.Printf("\nXML:\n%s", xmldoc)

	for i, element := range xmlReader.FindElements(nil, "Catalog", "Book", "Author") {
		fmt.Printf("\n/Catalog/Book/Author[%d] = %v", i, xmlReader.GetText(element, true))
	}

    //Genre won't get print here, as we haven't read it from XML file
	for i, element := range xmlReader.FindElements(nil, "Catalog", "Book", "Genre") {   
		fmt.Printf("\n/Catalog/Book/Genre[%d] = %v", i, xmlReader.GetText(element, true))
	}

XML:
<?xml version='1.0'?>
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
	
/Catalog/Book/Author[0] = Garghentini, Davide
/Catalog/Book/Author[1] = Garcia, Debra

```

## Updating XML File
```
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

Output:
Original XML:
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
</a>

Updated XML:
<a>
    <new_b>new_b_data</new_b>
    <c>c-data</c>
    <d>
        <e>e-data</e>
		<c>c-data</c>
    </d>
	<f>
		<g1>prepend_data</g1><g>g-data</g><g2>append_data</g2>
	</f>
</a>
```