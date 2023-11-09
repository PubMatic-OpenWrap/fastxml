package fastxml

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
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
		<g gk1="true" gk2="deleteme" gk3="10">g-data</g>
	</f>
</a>`)

	reader := NewXMLReader(nil)
	if err := reader.Parse(xmlDoc); err != nil {
		t.Errorf("xml parsing error: %v", err.Error())
		return
	}

	elementB := reader.FindElement(nil, "a", "b")
	elementF := reader.FindElement(nil, "a", "f")
	elementG := reader.FindElement(nil, "a", "f", "g")

	//xmlUpdater
	updater := NewXMLUpdater(xmlDoc)

	//remove elements
	updater.RemoveElement(reader.FindElement(nil, "a", "c"))

	//replace full element
	updater.ReplaceElement(elementB, `<new_b>new_b_data</new_b>`)

	//append or prepend new xml tag
	updater.PrependElement(elementF, `<f1>prepend_data</f1>`)
	updater.AppendElement(elementF, `<f2>append_data</f2>`)

	//append or prepend new xml tag in existing which has text
	updater.PrependElement(elementG, `<g1>prepend_tag</g1>`)
	updater.AppendElement(elementG, `<g2>append_tag</g2>`)

	//update text
	updater.UpdateText(elementG, "new-g-data")

	//add new attribute
	updater.AddAttribute(elementF, "", "fk1", "fv1")

	//update attribute name and value
	gk1 := reader.GetAttribute(elementG, "gk1")
	updater.UpdateAttributeName(gk1, "gk11")
	updater.UpdateAttributeValue(gk1, "false")

	//remove attribute
	gk2 := reader.GetAttribute(elementG, "gk2")
	updater.RemoveAttribute(gk2)

	//Build Updated XML File
	buf := bytes.Buffer{}
	updater.Build(&buf)

	t.Logf("\nOriginal XML:%s", xmlDoc)
	t.Logf("\n\nUpdated XML:%s", buf.String())
}

func TestXMLUpdater_AppendElement(t *testing.T) {
	type args struct {
		in         string
		operations func(xu *XMLUpdater, in []byte)
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "cannot_insert_in_empty_tag",
			args: args{
				in:         ``,
				operations: func(xu *XMLUpdater, in []byte) {},
			},
			want: ``,
		},
		{
			name: "append_inline_tag",
			args: args{
				in: `<a></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.AppendElement(reader.FindElement(nil, "a"), `<empty_tag/>`)
				},
			},
			want: `<a><empty_tag/></a>`,
		},
		{
			name: "append_tag",
			args: args{
				in: `<a></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.AppendElement(reader.FindElement(nil, "a"), `<tag>tagdata</tag>`)
				},
			},
			want: `<a><tag>tagdata</tag></a>`,
		},
		{
			name: "empty_element_tag",
			args: args{
				in: `<a></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					xu.AppendElement(nil, `<tag>tagdata</tag>`)
				},
			},
			want: `<a></a>`,
		},
		{
			name: "tag_with_text",
			args: args{
				in: `<a>test_data</a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.AppendElement(reader.FindElement(nil, "a"), `<tag>tagdata</tag>`)
				},
			},
			want: `<a>test_data<tag>tagdata</tag></a>`,
		},
		{
			name: "nested_tag_1",
			args: args{
				in: `<a><b><c>cdata</c></b></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.AppendElement(reader.FindElement(nil, "a", "b", "c"), `<tag>tagdata</tag>`)
				},
			},
			want: `<a><b><c>cdata<tag>tagdata</tag></c></b></a>`,
		},
		{
			name: "nested_tag_2",
			args: args{
				in: `<a><b><c>cdata</c></b></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.AppendElement(reader.FindElement(nil, "a", "b"), `<tag>tagdata</tag>`)
				},
			},
			want: `<a><b><c>cdata</c><tag>tagdata</tag></b></a>`,
		},
		{
			name: "multiple_elements",
			args: args{
				in: `<a><b>one</b></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					elementA := reader.FindElement(nil, "a")
					xu.AppendElement(elementA, `<b>two</b>`)
					xu.AppendElement(elementA, `<b>three</b>`)
				},
			},
			want: `<a><b>one</b><b>two</b><b>three</b></a>`,
		},
		{
			name: "multiple_nested_elements",
			args: args{
				in: `<a><b><c>cdata</c></b></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					elementB := reader.FindElement(nil, "a", "b")
					elementC := reader.FindElement(nil, "a", "b", "c")
					xu.AppendElement(elementB, `<b1>b1_data</b1>`)
					xu.AppendElement(elementC, `<c1>c1_data</c1>`)
				},
			},
			want: `<a><b><c>cdata<c1>c1_data</c1></c><b1>b1_data</b1></b></a>`,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xu := NewXMLUpdater([]byte(tt.args.in))
			tt.args.operations(xu, []byte(tt.args.in))
			//rebuild buffer
			out := bytes.Buffer{}
			xu.Build(&out)
			assert.Equal(t, tt.want, out.String())
		})
	}
}

func TestXMLUpdater_PrependElement(t *testing.T) {
	type args struct {
		in         string
		operations func(xu *XMLUpdater, in []byte)
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "cannot_insert_in_empty_tag",
			args: args{
				in:         ``,
				operations: func(xu *XMLUpdater, in []byte) {},
			},
			want: ``,
		},
		/*
			{
				name: "prepend_inline_tag",
				args: args{
					in: `<a ak1="av1"/>`,
					operations: func(xu *XMLUpdater, in []byte) {
						reader := NewXMLReader(nil)
						_ = reader.Parse(in)
						xu.PrependElement(reader.FindElement(nil, "a"), `<tag>tagdata</tag>`)
					},
				},
				want: `<a ak1="av1"><tag>tagdata</tag></a>`,
			},
		*/
		{
			name: "prepend_tag",
			args: args{
				in: `<a></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.PrependElement(reader.FindElement(nil, "a"), `<tag>tagdata</tag>`)
				},
			},
			want: `<a><tag>tagdata</tag></a>`,
		},
		{
			name: "empty_element_tag",
			args: args{
				in: `<a></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					xu.PrependElement(nil, `<tag>tagdata</tag>`)
				},
			},
			want: `<a></a>`,
		},
		{
			name: "tag_with_text",
			args: args{
				in: `<a>test_data</a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.PrependElement(reader.FindElement(nil, "a"), `<tag>tagdata</tag>`)
				},
			},
			want: `<a><tag>tagdata</tag>test_data</a>`,
		},
		{
			name: "nested_tag_1",
			args: args{
				in: `<a><b><c>cdata</c></b></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.PrependElement(reader.FindElement(nil, "a", "b", "c"), `<tag>tagdata</tag>`)
				},
			},
			want: `<a><b><c><tag>tagdata</tag>cdata</c></b></a>`,
		},
		{
			name: "nested_tag_2",
			args: args{
				in: `<a><b><c>cdata</c></b></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.PrependElement(reader.FindElement(nil, "a", "b"), `<tag>tagdata</tag>`)
				},
			},
			want: `<a><b><tag>tagdata</tag><c>cdata</c></b></a>`,
		},
		{
			name: "multiple_elements",
			args: args{
				in: `<a><b>one</b></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					elementA := reader.FindElement(nil, "a")
					xu.PrependElement(elementA, `<b>two</b>`)
					xu.PrependElement(elementA, `<b>three</b>`)
				},
			},
			want: `<a><b>two</b><b>three</b><b>one</b></a>`,
		},
		{
			name: "multiple_nested_elements",
			args: args{
				in: `<a><b><c>cdata</c></b></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					elementB := reader.FindElement(nil, "a", "b")
					elementC := reader.FindElement(nil, "a", "b", "c")
					xu.PrependElement(elementB, `<b1>b1_data</b1>`)
					xu.PrependElement(elementC, `<c1>c1_data</c1>`)
				},
			},
			want: `<a><b><b1>b1_data</b1><c><c1>c1_data</c1>cdata</c></b></a>`,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xu := NewXMLUpdater([]byte(tt.args.in))
			tt.args.operations(xu, []byte(tt.args.in))
			//rebuild buffer
			out := bytes.Buffer{}
			xu.Build(&out)
			assert.Equal(t, tt.want, out.String())
		})
	}
}

func TestXMLUpdater_ReplaceElement(t *testing.T) {
	type args struct {
		in         string
		operations func(xu *XMLUpdater, in []byte)
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "cannot_replace_in_empty_tag",
			args: args{
				in:         ``,
				operations: func(xu *XMLUpdater, in []byte) {},
			},
			want: ``,
		},
		{
			name: "replace_tag",
			args: args{
				in: `<a></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.ReplaceElement(reader.FindElement(nil, "a"), `<tag>tagdata</tag>`)
				},
			},
			want: `<tag>tagdata</tag>`,
		},
		{
			name: "replace_inline_tag",
			args: args{
				in: `<a/>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.ReplaceElement(reader.FindElement(nil, "a"), `<tag>tagdata</tag>`)
				},
			},
			want: `<tag>tagdata</tag>`,
		},
		{
			name: "empty_element_tag",
			args: args{
				in: `<a></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					xu.PrependElement(nil, `<tag>tagdata</tag>`)
				},
			},
			want: `<a></a>`,
		},
		{
			name: "tag_with_text",
			args: args{
				in: `<a ak1="av1">test_data</a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.ReplaceElement(reader.FindElement(nil, "a"), `<tag>tagdata</tag>`)
				},
			},
			want: `<tag>tagdata</tag>`,
		},
		{
			name: "nested_tag_1",
			args: args{
				in: `<a><b><c>cdata</c></b></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.ReplaceElement(reader.FindElement(nil, "a", "b", "c"), `<tag>tagdata</tag>`)
				},
			},
			want: `<a><b><tag>tagdata</tag></b></a>`,
		},
		{
			name: "nested_tag_2",
			args: args{
				in: `<a><b><c>cdata</c></b></a>`,
				operations: func(xu *XMLUpdater, in []byte) {
					reader := NewXMLReader(nil)
					_ = reader.Parse(in)
					xu.ReplaceElement(reader.FindElement(nil, "a", "b"), `<tag>tagdata</tag>`)
				},
			},
			want: `<a><tag>tagdata</tag></a>`,
		},
		/*
			{
				name: "invalid_replace_one_element_multiple_times",
				args: args{
					in: `<a><b>one</b></a>`,
					operations: func(xu *XMLUpdater, in []byte) {
						reader := NewXMLReader(nil)
						_ = reader.Parse(in)
						elementA := reader.FindElement(nil, "a")
						xu.ReplaceElement(elementA, `<b>two</b>`)
						xu.ReplaceElement(elementA, `<b>three</b>`)
					},
				},
				want: ``,
			},
			{
				name: "invalid_replace_overlapping_elements",
				args: args{
					in: `<a><b><c>cdata</c><d>ddata</d></b></a>`,
					operations: func(xu *XMLUpdater, in []byte) {
						reader := NewXMLReader(nil)
						_ = reader.Parse(in)
						elementB := reader.FindElement(nil, "a", "b")
						elementC := reader.FindElement(nil, "a", "b", "c")
						xu.ReplaceElement(elementB, `<b1>b1_data</b1>`)
						xu.ReplaceElement(elementC, `<c1>c1_data</c1>`)
					},
				},
				want: ``,
			},
		*/
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xu := NewXMLUpdater([]byte(tt.args.in))
			tt.args.operations(xu, []byte(tt.args.in))
			//rebuild buffer
			out := bytes.Buffer{}
			xu.Build(&out)
			assert.Equal(t, tt.want, out.String())
		})
	}
}
