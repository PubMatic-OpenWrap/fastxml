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

func TestXMLTokenizer(t *testing.T) {
	parser := XMLTokenizer{}
	in := []byte(minixml)
	tokenHandler := mockTokenHandler{}

	//parsing
	parser.Parse(in[:], tokenHandler.append)

	actual := getXML(in[:], tokenHandler.tokens[:])
	t.Logf("Raw Tags: \n%v\n", printTokens(in[:], tokenHandler.tokens[:]))
	t.Logf("XML: %v\n", getXML(in[:], tokenHandler.tokens[:]))
	assert.Equal(t, string(in), actual)
}

func TestXPathXMLTokenizer(t *testing.T) {
	parser := XMLTokenizer{}
	in := []byte(minixml)
	tokenHandler := mockTokenHandler{}

	//parsing
	parser.ParseWithXPath(in[:], xpaths["mini"], tokenHandler.append)

	actual := getXML(in[:], tokenHandler.tokens[:])
	t.Logf("Raw Tags: \n%v\n", printTokens(in[:], tokenHandler.tokens[:]))
	t.Logf("XML: %v\n", getXML(in[:], tokenHandler.tokens[:]))
	assert.Equal(t, string(in), actual)
}
