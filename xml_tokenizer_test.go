package fastxml

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	xml = `<VAST  version="3.0" xmlns:xs="http://www.w3.org/2001/XMLSchema">
    <Ad  id="20001">
        <Wrapper>
            <Error>http://example.com/error</Error>
            <Impression  id="Impression-ID">http://example.com/track/impression</Impression>
            <Creatives>
                <Creative  id="5480" sequence="1">
                    <Linear>
                        <Duration>00:00:16</Duration>
                        <TrackingEvents>
                            <Tracking  event="start">http://example.com/tracking/start</Tracking>
                            <Tracking  event="firstQuartile">http://example.com/tracking/firstQuartile</Tracking>
                            <Tracking  event="midpoint">http://example.com/tracking/midpoint</Tracking>
                            <Tracking  event="thirdQuartile">http://example.com/tracking/thirdQuartile</Tracking>
                            <Tracking  event="complete">http://example.com/tracking/complete</Tracking>
                            <Tracking  event="progress" offset="00:00:10">http://example.com/tracking/progress-10</Tracking>
                        </TrackingEvents>
                        <VideoClicks>
                            <ClickThrough  id="blog">
                                <![CDATA[https://iabtechlab.com]]>
                            </ClickThrough>
                        </VideoClicks>
                        <MediaFiles>
                            <MediaFile  id="5241" delivery="progressive" type="video/mp4" bitrate="500" width="400" height="300" minBitrate="360" maxBitrate="1080" scalable="1" maintainAspectRatio="1" codec="0" apiFramework="VAST">
                                <![CDATA[https://iab-publicfiles.s3.amazonaws.com/vast/VAST-4.0-Short-Intro.mp4]]>
                            </MediaFile>
                        </MediaFiles>
                    </Linear>
                </Creative>
            </Creatives>
        </Wrapper>
    </Ad>
</VAST>`
	minixml = `
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
</a>`
)

var xpaths = map[string]*xpath{
	"vast": GetXPath([][]string{
		{"VAST", "Ad", "InLine", "Impression"},
		{"VAST", "Ad", "InLine", "Error"},
		{"VAST", "Ad", "InLine", "Creatives", "Creative", "NonLinearAds", "TrackingEvents", "Tracking"},
		{"VAST", "Ad", "Wrapper", "Impression"},
		{"VAST", "Ad", "Wrapper", "Error"},
		{"VAST", "Ad", "Wrapper", "Creatives", "Creative", "Linear", "TrackingEvents", "Tracking"},
		{"VAST", "Ad", "Wrapper", "Creatives", "Creative", "Linear", "VideoClicks"},
	}),
	"mini": GetXPath([][]string{
		{"a", "b"},
		{"a", "d", "e"},
		{"a", "f"},
	}),
}

type mockTokenHandler struct {
	tokens []XMLToken
}

func (r *mockTokenHandler) append(_ string, parent *Element, child Element) {
	r.tokens = append(r.tokens, child.data)
}

func printTokens(in []byte, tokens []XMLToken) string {
	out := bytes.Buffer{}
	for i, token := range tokens {
		out.WriteString(fmt.Sprintf("%d:%s:end(%d:%d)\n", i, token.Name(in[:]), token.end.si, token.end.ei))
	}
	return out.String()
}

func getXML(in []byte, nodes []XMLToken) string {
	buf := bytes.Buffer{}
	start := 0
	for _, token := range nodes {
		buf.Write(in[start:token.end.ei])
		start = token.end.ei
	}
	return buf.String()
}

func TestXMLTokenizerParse(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantErr    bool
		expected   int      // expected number of tokens
		tokenNames []string // expected token names in order
	}{
		// Valid XML cases
		{
			name:       "simple_valid_xml",
			input:      "<root><child>data</child></root>",
			wantErr:    false,
			expected:   2,
			tokenNames: []string{"child", "root"},
		},
		{
			name:       "self_closing_tag",
			input:      "<root><child/></root>",
			wantErr:    false,
			expected:   2,
			tokenNames: []string{"child", "root"},
		},
		{
			name:       "empty_xml",
			input:      "",
			wantErr:    false,
			expected:   0,
			tokenNames: []string{},
		},
		{
			name:       "nested_elements",
			input:      "<a><b><c>data</c></b></a>",
			wantErr:    false,
			expected:   3,
			tokenNames: []string{"c", "b", "a"},
		},
		{
			name:       "with_attributes",
			input:      `<root id="1"><child type="text">data</child></root>`,
			wantErr:    false,
			expected:   2,
			tokenNames: []string{"child", "root"},
		},
		{
			name:       "without_handler",
			input:      "<root><child>data</child></root>",
			wantErr:    false,
			expected:   0,
			tokenNames: []string{},
		},
		{
			name:       "with_namespace",
			input:      "<ns:root><ns:child>data</ns:child></ns:root>",
			wantErr:    false,
			expected:   2,
			tokenNames: []string{"child", "root"},
		},
		// Invalid XML cases
		{
			name:       "missing_end_tag",
			input:      "<root><child></root>",
			wantErr:    true,
			expected:   0,
			tokenNames: []string{},
		},
		{
			name:       "mismatched_tags",
			input:      "<root><child></wrong></root>",
			wantErr:    true,
			expected:   0,
			tokenNames: []string{},
		},
		{
			name:       "incomplete_tag",
			input:      "<root><child",
			wantErr:    true,
			expected:   0,
			tokenNames: []string{},
		},
		// {
		// 	name:       "malformed_xml",
		// 	input:      "<<root>>",
		// 	wantErr:    true,
		// 	expected:   0,
		// 	tokenNames: []string{},
		// },
		{
			name:       "unclosed_tag",
			input:      "<root>",
			wantErr:    true,
			expected:   0,
			tokenNames: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := XMLTokenizer{}
			var handler *mockTokenHandler
			var callback TokenHandler

			if tt.name != "without_handler" {
				handler = &mockTokenHandler{}
				callback = handler.append
			}

			err := parser.Parse([]byte(tt.input), callback)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, errInvalidXML, err)
			} else {
				assert.NoError(t, err)
				if handler != nil {
					assert.Equal(t, tt.expected, len(handler.tokens), "unexpected number of tokens")

					// Verify token names
					actualNames := make([]string, len(handler.tokens))
					for i, token := range handler.tokens {
						actualNames[i] = string(token.Name([]byte(tt.input)))
					}
					assert.Equal(t, tt.tokenNames, actualNames, "token names mismatch")
				}
			}
		})
	}
}

func TestXMLTokenizer(t *testing.T) {
	xml := `<root><ns:child>data</ns:child></root>`
	parser := XMLTokenizer{}
	err := parser.Parse([]byte(xml), nil)
	assert.NoError(t, err)
}
